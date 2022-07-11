const NUMERIC = '0123456789';
const SMALL_ALPHABET = 'abcdefghijklmnopqrstuvwxyz';
const CAPITAL_ALPHABET = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

export const CAP_ONLY = NUMERIC + CAPITAL_ALPHABET;
export const ALPHANUMERIC = CAP_ONLY + SMALL_ALPHABET;

export const unsecure = (length: number, alphabet: string) => {
  let result = '';
  const charactersLength = alphabet.length;
  for (var i = 0; i < length; i++) {
    result += alphabet.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
};
