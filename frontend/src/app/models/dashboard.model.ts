export interface DashboardData {
  tenant: {
    name: string;
    companyCode: string;
    themeColor: string;
  };
  card: {
    maskedCardNumber: string;
    creditLimit: number;
    availableBalance: number;
    currency: string;
  };
  transactions: Array<{
    date: string;
    merchant: string;
    amount: number;
    status: string;
  }>;
}
