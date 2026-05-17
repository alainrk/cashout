// Chart rendering factories built on Chart.js.
(function (global) {
    const PALETTE = [
        '#4e79a7', '#f28e2b', '#e15759', '#76b7b2', '#59a14f',
        '#edc949', '#af7aa1', '#ff9da7', '#9c755f', '#bab0ab',
        '#8cd17d', '#b6992d', '#d37295', '#499894', '#a0cbe8',
        '#fabfd2', '#86bcb6', '#e19d9a', '#79706e',
    ];

    const categoryColorMap = new Map();
    function colorForCategory(cat) {
        if (!categoryColorMap.has(cat)) {
            categoryColorMap.set(cat, PALETTE[categoryColorMap.size % PALETTE.length]);
        }
        return categoryColorMap.get(cat);
    }

    function fmtEUR(value) {
        return new Intl.NumberFormat('en-US', {
            style: 'currency', currency: 'EUR', minimumFractionDigits: 2,
        }).format(value);
    }

    function destroy(canvas) {
        const existing = global.Chart && global.Chart.getChart(canvas);
        if (existing) existing.destroy();
    }

    function renderEmpty(canvas, message) {
        destroy(canvas);
        const ctx = canvas.getContext('2d');
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = '#888';
        ctx.font = '14px system-ui, sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(message || 'No data', canvas.width / 2, canvas.height / 2);
    }

    function renderCategoryDonut(canvas, byCategory, opts) {
        opts = opts || {};
        destroy(canvas);
        if (!byCategory || byCategory.length === 0) {
            renderEmpty(canvas, 'No data for this month');
            return null;
        }
        const labels = byCategory.map((c) => c.category);
        const data = byCategory.map((c) => c.amount);
        const colors = labels.map(colorForCategory);

        return new global.Chart(canvas, {
            type: 'doughnut',
            data: { labels, datasets: [{ data, backgroundColor: colors, borderWidth: 1 }] },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { position: 'right', labels: { boxWidth: 12 } },
                    tooltip: {
                        callbacks: {
                            label(ctx) {
                                const entry = byCategory[ctx.dataIndex];
                                return `${entry.category} · ${fmtEUR(entry.amount)} · ${entry.pct.toFixed(1)}%`;
                            },
                        },
                    },
                },
                onClick(evt, elements) {
                    if (!opts.onSliceClick || !elements.length) return;
                    const idx = elements[0].index;
                    opts.onSliceClick(byCategory[idx].category);
                },
            },
        });
    }

    function renderIncomeExpenseBar(canvas, income, expense) {
        destroy(canvas);
        if (!income && !expense) {
            renderEmpty(canvas, 'No data');
            return null;
        }
        return new global.Chart(canvas, {
            type: 'bar',
            data: {
                labels: ['Income', 'Expenses'],
                datasets: [{
                    data: [income, expense],
                    backgroundColor: ['#2ecc71', '#e74c3c'],
                    borderWidth: 0,
                }],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                indexAxis: 'y',
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        callbacks: { label: (ctx) => fmtEUR(ctx.parsed.x) },
                    },
                },
                scales: { x: { ticks: { callback: (v) => fmtEUR(v) } } },
            },
        });
    }

    function renderTrendLine(canvas, points, opts) {
        opts = opts || {};
        destroy(canvas);
        if (!points || points.length === 0) {
            renderEmpty(canvas, 'No data');
            return null;
        }
        const labels = points.map((p) => p.month);
        const income = points.map((p) => p.income);
        const expense = points.map((p) => p.expense);
        const balance = points.map((p) => p.balance);

        return new global.Chart(canvas, {
            type: 'line',
            data: {
                labels,
                datasets: [
                    { label: 'Income', data: income, borderColor: '#2ecc71', backgroundColor: '#2ecc71', tension: 0.2 },
                    { label: 'Expenses', data: expense, borderColor: '#e74c3c', backgroundColor: '#e74c3c', tension: 0.2 },
                    { label: 'Balance', data: balance, borderColor: '#3498db', backgroundColor: '#3498db', tension: 0.2 },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                interaction: { mode: 'index', intersect: false },
                plugins: {
                    tooltip: {
                        callbacks: { label: (ctx) => `${ctx.dataset.label}: ${fmtEUR(ctx.parsed.y)}` },
                    },
                },
                scales: { y: { ticks: { callback: (v) => fmtEUR(v) } } },
                onClick(evt, elements) {
                    if (!opts.onPointClick || !elements.length) return;
                    opts.onPointClick(points[elements[0].index].month);
                },
            },
        });
    }

    function renderYearStacked(canvas, byMonth) {
        destroy(canvas);
        if (!byMonth || byMonth.length === 0) {
            renderEmpty(canvas, 'No data for this year');
            return null;
        }
        const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
        const labels = byMonth.map((m) => monthNames[(m.month - 1) % 12]);
        const income = byMonth.map((m) => m.income);
        const expense = byMonth.map((m) => m.expense);
        return new global.Chart(canvas, {
            type: 'bar',
            data: {
                labels,
                datasets: [
                    { label: 'Income', data: income, backgroundColor: '#2ecc71' },
                    { label: 'Expenses', data: expense, backgroundColor: '#e74c3c' },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: { x: { stacked: false }, y: { stacked: false, ticks: { callback: (v) => fmtEUR(v) } } },
                plugins: {
                    tooltip: {
                        callbacks: { label: (ctx) => `${ctx.dataset.label}: ${fmtEUR(ctx.parsed.y)}` },
                    },
                },
            },
        });
    }

    global.CashoutCharts = {
        renderCategoryDonut,
        renderIncomeExpenseBar,
        renderTrendLine,
        renderYearStacked,
        colorForCategory,
    };
})(window);
