export interface Tenant {
  companyCode: string;
  name: string;
  features?: {
    paymentsEnabled?: boolean;
    profileEnabled?: boolean;
    [feature: string]: boolean | undefined;
  };
  theme: {
    primaryColor: string;
    secondaryColor: string;
    logoUrl: string;
  };
  contactAddress: string;
  mobileNumber: string;
  termsAndConditions: string;
  cardArt: {
    frontGradient: string;
    backGradient: string;
  };
  image: string;
}
