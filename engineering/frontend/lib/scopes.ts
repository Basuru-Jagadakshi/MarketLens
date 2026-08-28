export const appScopes = [
    
  // base
  'openid',
  'profile',

  // employment sectors
  'employment-sectors:create',
  'employment-sectors:update',
  'employment-sectors:delete',

  // experience
  'experiences:create',
  'experiences:update',
  'experiences:delete',

  // education level
  'education-levels:create',
  'education-levels:update',
  'education-levels:delete',

  // formality
  'formalities:create',
  'formalities:update',
  'formalities:delete',

  // gender
  'genders:create',
  'genders:update',
  'genders:delete',

  // vocational education
  'vocational-educations:create',
  'vocational-educations:update',
  'vocational-educations:delete',

  // major groups
  'major-groups:create',
  'major-groups:update',
  'major-groups:delete',

  // sub major groups
  'sub-major-groups:create',
  'sub-major-groups:update',
  'sub-major-groups:delete',

  // minor group
  'minor-groups:create',
  'minor-groups:update',
  'minor-groups:delete',

  // unit groups
  'unit-groups:create',
  'unit-groups:update',
  'unit-groups:delete',

  // occupation group
  'occupation-groups:create',
  'occupation-groups:update',
  'occupation-groups:delete',

  // industry sector
  'industry-sectors:create',
  'industry-sectors:update',
  'industry-sectors:delete',

  // industry division
  'industry-divisions:create',
  'industry-divisions:update',
  'industry-divisions:delete',

  // industry group
  'industry-groups:create',
  'industry-groups:update',
  'industry-groups:delete',

  // industry class
  'industry-classes:create',
  'industry-classes:update',
  'industry-classes:delete',

  // industry sub class
  'industry-sub-classes:create',
  'industry-sub-classes:update',
  'industry-sub-classes:delete',

] as const

export const scopeString = appScopes.join(' ')