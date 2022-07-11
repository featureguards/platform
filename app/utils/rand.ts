const NUMERIC = '0123456789';
const SMALL_ALPHABET = 'abcdefghijklmnopqrstuvwxyz';
const CAPITAL_ALPHABET = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

export const CAP_ONLY = CAPITAL_ALPHABET + NUMERIC;

export const unsecure = (length: number, alphabet: string) => {
  let result = '';
  const charactersLength = alphabet.length;
  for (var i = 0; i < length; i++) {
    result += alphabet.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
};
