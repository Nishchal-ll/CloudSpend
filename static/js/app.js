// CloudSpend — Client Application Logic

document.addEventListener('DOMContentLoaded', () => {
    // Set default date to today for new expenses
    const dateInput = document.getElementById('expenseDate');
    if (dateInput && !dateInput.value) {
        const today = new Date().toISOString().split('T')[0];
        dateInput.value = today;
    }

    // Close modals on Escape key
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            closeExpenseModal();
            closeBudgetModal();
            closeDeleteModal();
        }
    });

    // Close modal on backdrop click
    document.querySelectorAll('.modal-overlay, .modal-backdrop').forEach(modal => {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.classList.add('hidden');
            }
        });
    });
});

/* ==========================================================================
   Modal Functions: Add & Edit Expense
   ========================================================================== */

function openExpenseModal() {
    const modal = document.getElementById('expenseModal');
    const form = document.getElementById('expenseForm');
    const modalTitle = document.getElementById('modalTitle');
    const expenseId = document.getElementById('expenseId');
    const btnSubmit = document.getElementById('btnSubmitExpense');

    if (!modal) return;

    form.reset();
    expenseId.value = '';
    modalTitle.textContent = 'Add Cloud Expense';
    btnSubmit.textContent = 'Save Expense';

    const dateInput = document.getElementById('expenseDate');
    if (dateInput) {
        dateInput.value = new Date().toISOString().split('T')[0];
    }

    modal.classList.remove('hidden');
    document.getElementById('expenseTitle').focus();
}

function openEditModal(id, title, provider, category, amount, date, description) {
    const modal = document.getElementById('expenseModal');
    const modalTitle = document.getElementById('modalTitle');
    const btnSubmit = document.getElementById('btnSubmitExpense');

    if (!modal) return;

    document.getElementById('expenseId').value = id;
    document.getElementById('expenseTitle').value = title;
    document.getElementById('expenseProvider').value = provider || 'Other';
    document.getElementById('expenseCategory').value = category || 'Compute';
    document.getElementById('expenseAmount').value = amount;
    document.getElementById('expenseDate').value = date;
    document.getElementById('expenseDesc').value = description || '';

    modalTitle.textContent = 'Edit Cloud Expense #' + id;
    btnSubmit.textContent = 'Update Expense';

    modal.classList.remove('hidden');
}

function closeExpenseModal() {
    const modal = document.getElementById('expenseModal');
    if (modal) modal.classList.add('hidden');
}

async function handleExpenseSubmit(event) {
    event.preventDefault();

    const id = document.getElementById('expenseId').value;
    const title = document.getElementById('expenseTitle').value.trim();
    const provider = document.getElementById('expenseProvider').value;
    const category = document.getElementById('expenseCategory').value;
    const amount = parseFloat(document.getElementById('expenseAmount').value);
    const date = document.getElementById('expenseDate').value;
    const description = document.getElementById('expenseDesc').value.trim();

    if (!title || isNaN(amount) || amount <= 0 || !date) {
        showToast('Please fill out all required fields with valid values.', 'error');
        return;
    }

    const payload = {
        title: title,
        provider: provider,
        category: category,
        amount: amount,
        expense_date: date,
        description: description
    };

    const isEdit = Boolean(id);
    const url = isEdit ? `/api/expenses/${id}` : '/api/expenses';
    const method = isEdit ? 'PUT' : 'POST';

    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to save expense');
        }

        showToast(isEdit ? 'Expense updated successfully!' : 'Expense added successfully!', 'success');
        closeExpenseModal();

        // Reload page after short delay to reflect changes in SSR tables and KPI cards
        setTimeout(() => {
            window.location.reload();
        }, 600);
    } catch (err) {
        showToast(err.message, 'error');
    }
}

/* ==========================================================================
   Modal Functions: Delete Expense
   ========================================================================== */

function confirmDeleteExpense(id, title) {
    const modal = document.getElementById('deleteModal');
    if (!modal) return;

    document.getElementById('deleteExpenseId').value = id;
    document.getElementById('deleteExpenseName').textContent = `"${title}" (#${id})`;
    modal.classList.remove('hidden');
}

function closeDeleteModal() {
    const modal = document.getElementById('deleteModal');
    if (modal) modal.classList.add('hidden');
}

async function executeDeleteExpense() {
    const id = document.getElementById('deleteExpenseId').value;
    if (!id) return;

    try {
        const response = await fetch(`/api/expenses/${id}`, {
            method: 'DELETE'
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || 'Failed to delete expense');
        }

        showToast('Expense deleted successfully!', 'success');
        closeDeleteModal();

        // Remove row immediately or reload
        const row = document.getElementById(`expense-row-${id}`);
        if (row) {
            row.style.opacity = '0';
            setTimeout(() => {
                window.location.reload();
            }, 500);
        } else {
            setTimeout(() => {
                window.location.reload();
            }, 600);
        }
    } catch (err) {
        showToast(err.message, 'error');
    }
}

/* ==========================================================================
   Modal Functions: Adjust Budget
   ========================================================================== */

function openBudgetModal(currentBudget) {
    const modal = document.getElementById('budgetModal');
    if (!modal) return;

    const input = document.getElementById('budgetInput');
    if (input && currentBudget) {
        input.value = parseFloat(currentBudget).toFixed(0);
    }
    modal.classList.remove('hidden');
}

function closeBudgetModal() {
    const modal = document.getElementById('budgetModal');
    if (modal) modal.classList.add('hidden');
}

async function handleBudgetSubmit(event) {
    event.preventDefault();
    const budgetInput = document.getElementById('budgetInput');
    const newLimit = parseFloat(budgetInput.value);

    if (isNaN(newLimit) || newLimit <= 0) {
        showToast('Please enter a valid positive budget amount.', 'error');
        return;
    }

    try {
        const response = await fetch('/api/budget', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ monthly_limit: newLimit })
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || 'Failed to update budget limit');
        }

        showToast('Monthly budget updated!', 'success');
        closeBudgetModal();
        setTimeout(() => {
            window.location.reload();
        }, 600);
    } catch (err) {
        showToast(err.message, 'error');
    }
}

/* ==========================================================================
   Toast Notification System
   ========================================================================== */

function showToast(message, type = 'success') {
    const container = document.getElementById('toastContainer');
    if (!container) return;

    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    
    const icon = type === 'success' 
        ? '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2"><path d="M20 6 9 17l-5-5"/></svg>'
        : '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#f43f5e" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" x2="12" y1="8" y2="12"/><line x1="12" x2="12.01" y1="16" y2="16"/></svg>';

    toast.innerHTML = `${icon}<span>${message}</span>`;
    container.appendChild(toast);

    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transform = 'translateY(10px)';
        toast.style.transition = 'all 0.3s ease';
        setTimeout(() => toast.remove(), 300);
    }, 3500);
}
