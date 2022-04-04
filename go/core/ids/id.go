package ids

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/spaolacci/murmur3"
)

type ObjectType int
type ID string
type Shard int
type PhysicalShard int

const (
	// Never add in the middle.
	Unknown ObjectType = iota
	User
	Project
	ProjectMember
	ProjectInvite
	Environment
	FeatureToggle

	// Add here
	MaxObjectType
)

const (
	// Bits
	TotalBits = 78
	// LogicalShard + Reserved = 2 Chars
	LogicalShardBits = 9
	ReservedBits     = 3
	TypeBits         = 12 // 2 CHARs
	NonRandomBits    = LogicalShardBits + ReservedBits + TypeBits
	RandomBits       = TotalBits - NonRandomBits

	// Shard limit
	MaxLogicalShards = 1 << LogicalShardBits
	MaxRand          = 1 << RandomBits
)

var (
	// use a custom base-64 for efficiency in operations (base 2)
	alphabet                    = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@%")
	base                        = uint64(len(alphabet))
	alphabetMap map[rune]uint64 = make(map[rune]uint64, len(alphabet))

	IDLen = int(TotalBits / base) // Chars (Not Bits)
)

func (s Shard) validate() error {
	if s >= MaxLogicalShards {
		return errors.WithStack(fmt.Errorf("%s exceeds maximum allowed logical shards", s))
	}
	if s < 0 {
		return errors.WithStack(fmt.Errorf("%s must be >= 0", s))
	}
	return nil
}

func (id ID) String() string {
	return string(id)
}

func (s Shard) String() string {
	return fmt.Sprintf("shard-%d", s)
}

func (o ObjectType) Validate() error {
	if int(o) <= int(Unknown) || int(o) >= int(MaxObjectType) {
		return errors.WithStack(fmt.Errorf("unexpected object type: %s", o))
	}
	return nil
}

func (o ObjectType) String() string {
	switch o {
	case Project:
		return "Project"
	case ProjectInvite:
		return "ProjectInvite"
	case Environment:
		return "Environment"
	case FeatureToggle:
		return "FeatureToggle"
	case User:
		return "User"
	default:
		return "Unknown"
	}
}

func (o ObjectType) IsRoot() bool {
	return o == Project
}

type IDs struct {
	physicalShards int
}

type IDsOpts struct {
	PhysicalBits int

	// Testing
	Seed int64
}

func New(opts IDsOpts) (*IDs, error) {
	seed := opts.Seed
	if seed <= 0 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	if opts.PhysicalBits >= LogicalShardBits {
		return nil, errors.WithStack(fmt.Errorf("physical shards must be <= maximum logical shards"))
	}
	if opts.PhysicalBits <= 0 {
		return nil, errors.WithStack(fmt.Errorf("physical shards must be positive"))
	}
	return &IDs{
		physicalShards: 1 << opts.PhysicalBits,
	}, nil
}

// reverse reverses in-place
func reverse(s []rune) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func encodePrepended(n uint64, expected int) ([]rune, error) {
	encoded, err := encode(n)
	if err != nil {
		return nil, err
	}

	l := len(encoded)
	if l > expected {
		return nil, errors.WithStack(fmt.Errorf("%d was encoded to %s which was longer than %d", n, string(encoded), expected))
	}

	if l == expected {
		return encoded, nil
	}

	// Prepend
	prepended := make([]rune, expected-l, expected)
	for i := range prepended {
		prepended[i] = alphabet[0]
	}

	return append(prepended, encoded...), nil

}

// encode encodes a positive number into alphabet and return the string.
func encode(n uint64) ([]rune, error) {
	if n == 0 {
		return []rune{alphabet[0]}, nil
	}

	res := make([]rune, 0, IDLen)
	for n > 0 {
		rem := n % base
		n = n / base
		res = append(res, alphabet[rem])
	}
	reverse(res)
	return res, nil
}

func decode(s []rune) (uint64, error) {
	n := uint64(0)
	l := len(s)
	for i, c := range s {
		pow := (l - (i + 1))
		idx, ok := alphabetMap[c]
		if !ok {
			return 0, errors.WithStack(fmt.Errorf("unknown character %c not in alphabet", c))
		}
		n += idx * (base << pow)
	}
	return n, nil
}

func RandomID(t ObjectType) (ID, error) {
	if !t.IsRoot() {
		return "", errors.WithStack(fmt.Errorf("object must be root"))
	}
	shard := Shard(rand.Intn(MaxLogicalShards))
	return IDFromShard(shard, t)
}

func Parse(id ID) (ObjectType, Shard, error) {
	raw := []rune(id)
	if len(id) != IDLen {
		return Unknown, -1, errors.WithStack(fmt.Errorf("unexpected ID len for %s", id))
	}
	ot, err := decode(raw[:2])
	if err != nil {
		return Unknown, -1, err
	}
	if err := ObjectType(ot).Validate(); err != nil {
		return Unknown, -1, err
	}
	shardAndReserved, err := decode(raw[2:4])
	if err != nil {
		return Unknown, -1, err
	}
	shard := Shard(shardAndReserved & (MaxLogicalShards - 1))
	return ObjectType(ot), shard, nil
}

func IDFromRoot(rootID ID, ot ObjectType) (ID, error) {
	rootObjectType, shard, err := Parse(rootID)
	if err != nil {
		return "", err
	}
	if !rootObjectType.IsRoot() {
		return "", errors.WithStack(fmt.Errorf("given root %s is not a root", rootID))
	}
	return IDFromShard(shard, ot)
}

//   IDFromShard generates an ID of IDLen chars that has the the following format:
//       2-char            +                2-char            +    9-char
//   12-bits (object_type) | 3-bits (reserved) | 9-bits shard | 54-bits (random)
func IDFromShard(shard Shard, t ObjectType) (ID, error) {
	if err := shard.validate(); err != nil {
		return "", err
	}
	encodedType, err := encodePrepended(uint64(t), 2)
	if err != nil {
		return "", err
	}
	// Here is where we would add reserved. Right now, they're 0.
	encodedShard, err := encodePrepended(uint64(shard), 2)
	if err != nil {
		return "", err
	}
	encodedRand, err := encodePrepended(uint64(rand.Int63n(MaxRand)), 9)
	if err != nil {
		return "", err
	}
	return ID(append(append(encodedType, encodedShard...), encodedRand...)), nil
}

func (ids *IDs) PhysicalShardFromShard(shard Shard) (PhysicalShard, error) {
	if err := shard.validate(); err != nil {
		return -1, err
	}
	return PhysicalShard((ids.physicalShards - 1) & int(shard)), nil
}

func ShardIDFromKey(key string) Shard {
	return Shard(murmur3.Sum64([]byte(key)) % uint64(MaxLogicalShards))
}

func init() {
	for i, c := range alphabet {
		if _, ok := alphabetMap[c]; ok {
			log.Fatalf("multiple occurances of %s in alphabet", string(c))
		}
		alphabetMap[c] = uint64(i)
	}

	if RandomBits <= 0 {
		log.Fatal("random bits must be positive")
	}
}
