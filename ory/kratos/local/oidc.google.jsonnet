local claims = {
} + std.extVar('claims');

{
  identity: {
    traits: {
      [if "email" in claims && claims.email_verified then "email" else null]: claims.email,
      // additional claims
      email_verified: claims.email_verified,
      // please also see the `Google specific claims` section
      first_name: claims.given_name,
      last_name: claims.family_name,
      [if "hd" in claims && claims.email_verified then "hd" else null]: claims.hd,
      profile: claims.picture,
    },
  },
}