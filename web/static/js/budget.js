// Budget tab: monthly budget management.
(function () {
  const statusEl = document.getElementById('budgetStatus');
  const form = document.getElementById('budgetForm');
  const amountInput = document.getElementById('budgetAmount');
  const submitBtn = document.getElementById('submitBudgetBtn');
  const deleteBtn = document.getElementById('deleteBudgetBtn');
  const messageEl = document.getElementById('budgetMessage');

  if (!statusEl || !form) return;

  const fmtEUR = (n) =>
    new Intl.NumberFormat(undefined, {
      style: 'currency',
      currency: 'EUR',
    }).format(n);

  function showMessage(text, kind) {
    messageEl.textContent = text;
    messageEl.className = 'message ' + (kind || 'success');
    setTimeout(() => {
      messageEl.textContent = '';
      messageEl.className = 'message';
    }, 4000);
  }

  function renderStatus(data) {
    if (!data || !data.hasBudget) {
      statusEl.innerHTML =
        '<div class="budget-empty">No monthly budget set yet.</div>';
      deleteBtn.hidden = true;
      submitBtn.textContent = 'Save Budget';
      return;
    }

    const pct = Math.max(0, Math.min(100, data.pct));
    let state = 'ok';
    if (data.pct >= 100) state = 'over';
    else if (data.pct >= 80) state = 'warn';

    statusEl.innerHTML = `
      <div class="budget-card">
        <div class="budget-card-head">
          <div class="budget-amounts">
            <span class="budget-spent">${fmtEUR(data.spent)}</span>
            <span class="budget-of">of</span>
            <span class="budget-total">${fmtEUR(data.amount)}</span>
          </div>
          <div class="budget-meta">
            <span class="budget-pct budget-pct--${state}">${data.pct}%</span>
            <span class="budget-month">${data.month}</span>
          </div>
        </div>
        <div class="budget-bar">
          <div class="budget-bar-fill budget-bar-fill--${state}" style="width:${pct}%"></div>
        </div>
      </div>
    `;
    amountInput.value = data.amount.toString();
    deleteBtn.hidden = false;
    submitBtn.textContent = 'Update Budget';
  }

  async function fetchBudget() {
    try {
      const res = await fetch('/web/api/budget', { credentials: 'same-origin' });
      const json = await res.json();
      if (!res.ok) throw new Error(json.error || 'Request failed');
      renderStatus(json);
    } catch (e) {
      statusEl.innerHTML = '<div class="error">Failed to load budget.</div>';
    }
  }

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const raw = (amountInput.value || '').replace(',', '.').trim();
    const amount = parseFloat(raw);
    if (!(amount > 0)) {
      showMessage('Please enter a positive amount.', 'error');
      return;
    }
    submitBtn.disabled = true;
    try {
      const res = await fetch('/web/api/budget', {
        method: 'POST',
        credentials: 'same-origin',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ amount }),
      });
      const json = await res.json();
      if (!res.ok) throw new Error(json.error || 'Save failed');
      renderStatus(json);
      showMessage('Budget saved.', 'success');
    } catch (e) {
      showMessage('Failed to save budget.', 'error');
    } finally {
      submitBtn.disabled = false;
    }
  });

  deleteBtn.addEventListener('click', async () => {
    if (!confirm('Remove your monthly budget?')) return;
    deleteBtn.disabled = true;
    try {
      const res = await fetch('/web/api/budget', {
        method: 'DELETE',
        credentials: 'same-origin',
      });
      const json = await res.json();
      if (!res.ok) throw new Error(json.error || 'Delete failed');
      amountInput.value = '';
      renderStatus(json);
      showMessage('Budget removed.', 'success');
    } catch (e) {
      showMessage('Failed to remove budget.', 'error');
    } finally {
      deleteBtn.disabled = false;
    }
  });

  // Lazy-load on first tab activation, or immediately if Budget is the persisted current page.
  let loaded = false;
  function ensureLoaded() {
    if (loaded) return;
    loaded = true;
    fetchBudget();
  }
  document.querySelectorAll('.nav-tab').forEach((tab) => {
    tab.addEventListener('click', () => {
      if (tab.dataset.page === 'budget') ensureLoaded();
    });
  });
  if (
    (localStorage.getItem('currentPage') || 'transactions') === 'budget'
  ) {
    ensureLoaded();
  }
})();
