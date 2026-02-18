import { Tenant } from '../models/tenant.model';

export const DEFAULT_TENANT: Tenant = {
  companyCode: 'cardflex',
  name: 'CardFlex',
  theme: {
    primaryColor: '#00539C',
    secondaryColor: '#8FB4D8',
    logoUrl: 'https://placehold.co/132x36/00539C/ffffff?text=CardFlex'
  },
  contactAddress: '100 Main St, New York, NY 10001',
  mobileNumber: '+1 (800) 555-0100',
  termsAndConditions: 'Standard CardFlex terms apply. Statement generated monthly and payment due within 25 days.',
  cardArt: {
    frontGradient: 'linear-gradient(135deg, #fff8e7, #f8b500)',
    backGradient: 'linear-gradient(135deg, #c2e9fb, #81a4fd)'
  },
  image: 'https://images.unsplash.com/photo-1556740749-887f6717d7e4?auto=format&fit=crop&w=1200&q=80'
};

export const TENANT_CONFIGS: Record<string, Tenant> = {
  'chase-bank': {
    companyCode: 'chase-bank',
    name: 'Chase Bank',
    theme: {
      primaryColor: '#0A2A66',
      secondaryColor: '#2E8BC0',
      logoUrl: 'https://placehold.co/132x36/0A2A66/ffffff?text=Chase+Bank'
    },
    contactAddress: '270 Park Ave, New York, NY 10017',
    mobileNumber: '+1 (800) 935-9935',
    termsAndConditions: 'APR, fees, and rewards are subject to Chase Bank cardmember agreement and credit approval.',
    cardArt: {
      frontGradient: 'linear-gradient(135deg, #dbeafe, #60a5fa)',
      backGradient: 'linear-gradient(135deg, #0A2A66, #2E8BC0)'
    },
    image: 'https://images.unsplash.com/photo-1556745757-8d76bdb6984b?auto=format&fit=crop&w=1200&q=80'
  },
  'wells-fargo': {
    companyCode: 'wells-fargo',
    name: 'Wells Fargo',
    theme: {
      primaryColor: '#B31B1B',
      secondaryColor: '#F2C94C',
      logoUrl: 'https://placehold.co/132x36/B31B1B/ffffff?text=Wells+Fargo'
    },
    contactAddress: '420 Montgomery St, San Francisco, CA 94104',
    mobileNumber: '+1 (800) 869-3557',
    termsAndConditions: 'Wells Fargo credit products are subject to account agreement, eligibility, and applicable state laws.',
    cardArt: {
      frontGradient: 'linear-gradient(135deg, #fee2e2, #ef4444)',
      backGradient: 'linear-gradient(135deg, #B31B1B, #F2C94C)'
    },
    image: 'https://images.unsplash.com/photo-1565514020179-026b92b84bb6?auto=format&fit=crop&w=1200&q=80'
  },
  'capital-one': {
    companyCode: 'capital-one',
    name: 'Capital One',
    theme: {
      primaryColor: '#003B95',
      secondaryColor: '#D62728',
      logoUrl: 'https://placehold.co/132x36/003B95/ffffff?text=Capital+One'
    },
    contactAddress: '1680 Capital One Dr, McLean, VA 22102',
    mobileNumber: '+1 (877) 383-4802',
    termsAndConditions: 'Capital One card terms, credit limits, and rewards are governed by your signed card agreement.',
    cardArt: {
      frontGradient: 'linear-gradient(135deg, #dbeafe, #2563eb)',
      backGradient: 'linear-gradient(135deg, #003B95, #D62728)'
    },
    image: 'https://images.unsplash.com/photo-1601597111158-2fceff292cdc?auto=format&fit=crop&w=1200&q=80'
  }
};
