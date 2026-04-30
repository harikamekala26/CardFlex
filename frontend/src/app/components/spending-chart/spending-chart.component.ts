import { CommonModule, CurrencyPipe } from '@angular/common';
import { Component, Input } from '@angular/core';

import { DashboardTransaction } from '../../models/dashboard.model';

export interface SpendingSummaryItem {
  month: string;
  amount: number;
  percent: number;
}

@Component({
  selector: 'app-spending-chart',
  standalone: true,
  imports: [CommonModule, CurrencyPipe],
  templateUrl: './spending-chart.component.html',
  styleUrl: './spending-chart.component.css'
})
export class SpendingChartComponent {
  @Input() transactions: DashboardTransaction[] = [];
  @Input() currency = 'USD';

  get hasTransactions(): boolean {
    return this.spendingSummary.length > 0;
  }

  get spendingSummary(): SpendingSummaryItem[] {
    const monthlySpending = new Map<string, { date: Date; amount: number }>();

    for (const transaction of this.transactions) {
      const date = new Date(transaction.date);

      if (Number.isNaN(date.getTime())) {
        continue;
      }

      const monthKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
      const existing = monthlySpending.get(monthKey);
      const spendAmount = Math.abs(transaction.amount);

      monthlySpending.set(monthKey, {
        date: existing?.date ?? new Date(date.getFullYear(), date.getMonth(), 1),
        amount: (existing?.amount ?? 0) + spendAmount
      });
    }

    const items = Array.from(monthlySpending.values()).sort((first, second) => first.date.getTime() - second.date.getTime());
    const maxAmount = Math.max(...items.map((item) => item.amount), 0);

    return items.map((item) => ({
      month: item.date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' }),
      amount: item.amount,
      percent: maxAmount > 0 ? Math.round((item.amount / maxAmount) * 100) : 0
    }));
  }

  trackSpendingMonth(_: number, item: SpendingSummaryItem): string {
    return item.month;
  }
}
