export interface DashboardTenant {
  name: string;
  companyCode: string;
  themeColor: string;
}

export interface DashboardAccountSummary {
  maskedCardNumber: string;
  creditLimit: number;
  availableBalance: number;
  currency: string;
}

export interface DashboardTransaction {
  date: string;
  merchant: string;
  amount: number;
  status: string;
}

export interface DashboardApiResponse {
  tenant: DashboardTenant;
  accountSummary?: DashboardAccountSummary | null;
  card?: DashboardAccountSummary | null;
  transactions: DashboardTransaction[];
}

export interface DashboardData {
  tenant: DashboardTenant;
  accountSummary: DashboardAccountSummary;
  transactions: DashboardTransaction[];
}

export function normalizeDashboardData(response: DashboardApiResponse): DashboardData {
  return {
    tenant: response.tenant,
    accountSummary: response.accountSummary ?? response.card ?? createEmptyAccountSummary(),
    transactions: response.transactions ?? []
  };
}

function createEmptyAccountSummary(): DashboardAccountSummary {
  return {
    maskedCardNumber: '',
    creditLimit: 0,
    availableBalance: 0,
    currency: 'USD'
  };
}
