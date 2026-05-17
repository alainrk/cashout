// Data layer for analytics endpoints with in-memory cache.
(function (global) {
    const monthlyCache = new Map();
    const trendCache = new Map();
    const yearCache = new Map();

    async function fetchJSON(url) {
        const response = await fetch(url);
        const data = await response.json();
        if (!response.ok) throw new Error(data.error || 'Request failed');
        return data;
    }

    async function fetchMonthly(ym) {
        if (monthlyCache.has(ym)) return monthlyCache.get(ym);
        const data = await fetchJSON('/web/api/analytics/monthly?month=' + encodeURIComponent(ym));
        monthlyCache.set(ym, data);
        return data;
    }

    async function fetchTrend(months) {
        const key = String(months);
        if (trendCache.has(key)) return trendCache.get(key);
        const data = await fetchJSON('/web/api/analytics/trend?months=' + encodeURIComponent(months));
        trendCache.set(key, data);
        return data;
    }

    async function fetchYear(year) {
        const key = String(year);
        if (yearCache.has(key)) return yearCache.get(key);
        const data = await fetchJSON('/web/api/analytics/year?year=' + encodeURIComponent(year));
        yearCache.set(key, data);
        return data;
    }

    function invalidate() {
        monthlyCache.clear();
        trendCache.clear();
        yearCache.clear();
    }

    global.Analytics = { fetchMonthly, fetchTrend, fetchYear, invalidate };
})(window);
