document.addEventListener('DOMContentLoaded', () => {
    // Logs Table Element
    const logsTableBody = document.getElementById('logsTableBody');

    // --- LIVE USAGE LOGS POLLING ---
    async function fetchLogs() {
        try {
            const res = await fetch('/admin/logs');
            const logs = await res.json();
            renderLogs(logs);
        } catch (e) {
            console.error("Failed to fetch logs", e);
        }
    }

    function renderLogs(logs) {
        if (!logs || logs.length === 0) {
            logsTableBody.innerHTML = `<tr><td colspan="3" style="text-align:center; color: #64748b;">No usage logs available yet. Send some requests!</td></tr>`;
            return;
        }

        logsTableBody.innerHTML = logs.map(log => {
            const date = new Date(log.timestamp).toLocaleTimeString();
            const badgeClass = log.allowed ? 'badge-allowed' : 'badge-blocked';
            const badgeText = log.allowed ? 'Allowed' : 'Rate Limited';
            return `
                <tr>
                    <td style="color: #64748b;">${date}</td>
                    <td style="font-family: monospace; color: #38bdf8;">${log.identifier}</td>
                    <td><span class="badge ${badgeClass}">${badgeText}</span></td>
                </tr>
            `;
        }).join('');
    }

    // Poll every 2 seconds
    fetchLogs();
    setInterval(fetchLogs, 2000);

    // --- DYNAMIC HOST UPDATE ---
    // Update the code snippets to show the actual deployed URL
    const currentHost = window.location.origin;
    if (currentHost && !currentHost.includes("localhost") && !currentHost.includes("127.0.0.1")) {
        const codeBlocks = document.querySelectorAll('pre code');
        codeBlocks.forEach(block => {
            block.innerHTML = block.innerHTML.replace(/http:\/\/localhost:8080/g, currentHost);
        });
    }

    // --- TABS LOGIC ---
    const tabBtns = document.querySelectorAll('.tab-btn');

    if (tabBtns.length > 0) {
        tabBtns.forEach(btn => {
            btn.addEventListener('click', () => {
                const group = btn.dataset.group;
                
                // Get all buttons in the same group
                const groupBtns = document.querySelectorAll(`.tab-btn[data-group="${group}"]`);
                
                // Get all content targets for this group
                const targets = Array.from(groupBtns).map(b => b.dataset.target);
                
                // Remove active class from buttons and contents in this group
                groupBtns.forEach(b => b.classList.remove('active'));
                targets.forEach(id => {
                    const el = document.getElementById(id);
                    if(el) el.classList.remove('active');
                });

                // Add active class to clicked button and its target content
                btn.classList.add('active');
                const targetEl = document.getElementById(btn.dataset.target);
                if(targetEl) targetEl.classList.add('active');
            });
        });
    }
});
