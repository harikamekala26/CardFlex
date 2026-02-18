export interface Tenant {
  companyCode: string;
  name: string;
  primaryColor: string;
  secondaryColor: string;
  contactAddress: string;
  mobileNumber: string;
  termsAndConditions: string;
  cardArt: {
    frontGradient: string;
    backGradient: string;
  };
  image: string;
}
