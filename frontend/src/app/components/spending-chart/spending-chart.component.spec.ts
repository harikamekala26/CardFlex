import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';

import { SpendingChartComponent } from './spending-chart.component';

describe('SpendingChartComponent', () => {
  let fixture: ComponentFixture<SpendingChartComponent>;
  let component: SpendingChartComponent;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SpendingChartComponent]
    }).compileComponents();

    fixture = TestBed.createComponent(SpendingChartComponent);
    component = fixture.componentInstance;
  });

  it('renders an empty state when there are no chartable transactions', () => {
    component.transactions = [];

    fixture.detectChanges();

    expect(component.hasTransactions).toBeFalse();
    expect(component.spendingSummary).toEqual([]);
    expect(fixture.nativeElement.textContent).toContain('No transactions to display');
    expect(fixture.debugElement.queryAll(By.css('.spending-bar')).length).toBe(0);
  });

  it('groups transaction amounts by month and renders chart bars', () => {
    component.transactions = [
      { date: '2026-02-14', merchant: 'Grocery Mart', amount: -82.41, status: 'Posted' },
      { date: '2026-02-20', merchant: 'Coffee Stand', amount: -17.59, status: 'Posted' },
      { date: '2026-03-03', merchant: 'Payment', amount: 50, status: 'Posted' },
      { date: 'not-a-date', merchant: 'Ignored', amount: -999, status: 'Posted' }
    ];

    fixture.detectChanges();

    expect(component.spendingSummary).toEqual([
      { month: 'Feb 2026', amount: 100, percent: 100 },
      { month: 'Mar 2026', amount: 50, percent: 50 }
    ]);
    expect(fixture.nativeElement.textContent).toContain('Feb 2026');
    expect(fixture.nativeElement.textContent).toContain('Mar 2026');
    expect(fixture.nativeElement.textContent).toContain('$100');
    expect(fixture.debugElement.queryAll(By.css('.spending-bar')).length).toBe(2);
  });
});
