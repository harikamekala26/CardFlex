export interface Tenant {
  companyCode: string;
  name: string;
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
