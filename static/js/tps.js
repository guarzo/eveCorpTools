// Function to truncate labels
function truncateLabel(label, length) {
    if (!label || typeof label !== 'string') {
        return '';
    }
    return label.length > length ? label.substring(0, length) + '...' : label;
}

// Function to get colors for datasets
function getColor(index) {
    const colors = [
        '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0',
        '#9966FF', '#FF9F40', '#E7E9ED', '#76D7C4',
        '#C0392B', '#8E44AD', '#2ECC71', '#1ABC9C',
        '#3498DB', '#F1C40F', '#E67E22', '#95A5A6',
    ];
    return colors[index % colors.length];
}

// Wrap all code in DOMContentLoaded listener
document.addEventListener('DOMContentLoaded', function() {
    // Variable to keep track of current time frame
    let currentTimeFrame = 'mtd';

    // Set initial active button
    document.getElementById('btn-mtd').classList.add('active');

    // Function to set the time frame
    window.setTimeFrame = function(timeFrame) {
        currentTimeFrame = timeFrame;
        updateAllCharts();

        // Update button styles
        document.getElementById('btn-mtd').classList.remove('active');
        document.getElementById('btn-lastMonth').classList.remove('active');
        document.getElementById('btn-ytd').classList.remove('active');

        if (timeFrame === 'mtd') {
            document.getElementById('btn-mtd').classList.add('active');
        } else if (timeFrame === 'lastMonth') {
            document.getElementById('btn-lastMonth').classList.add('active');
        } else if (timeFrame === 'ytd') {
            document.getElementById('btn-ytd').classList.add('active');
        }
    };

    // Initialize all charts
    let damageFinalBlowsChart = null;
    let ourLossesCombinedChart = null;
    let characterPerformanceChart = null;
    let ourShipsUsedChart = null;
    let victimsSunburstChart = null;
    let killActivityChart = null;
    let killHeatmapChart = null;
    let killLossRatioChart = null;
    let topShipsKilledChart = null;
    let valueOverTimeChart = null;

    // Function to update all charts
    function updateAllCharts() {
        updateDamageFinalBlowsChart();
        updateOurLossesCombinedChart();
        updateCharacterPerformanceChart();
        updateOurShipsUsedChart();
        updateVictimsSunburstChart();
        updateKillActivityChart();
        updateKillHeatmapChart();
        updateKillLossRatioChart();
        updateTopShipsKilledChart();
        updateValueOverTimeChart();
    }

    // 1. Damage Done and Final Blows Chart
    function updateDamageFinalBlowsChart() {
        let data;
        if (currentTimeFrame === 'mtd') {
            data = window.mtdCharacterDamageData;
        } else if (currentTimeFrame === 'ytd') {
            data = window.ytdCharacterDamageData;
        } else if (currentTimeFrame === 'lastMonth') {
            data = window.lastMCharacterDamageData;
        }

        if (!data) return;

        const labels = data.map(item => item.Name);
        const damageData = data.map(item => item.DamageDone);
        const finalBlowsData = data.map(item => item.FinalBlows);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const ctxElem = document.getElementById('damageFinalBlowsChart');
        if (!ctxElem) return;

        if (damageFinalBlowsChart) {
            // Update existing chart
            damageFinalBlowsChart.data.labels = truncatedLabels;
            damageFinalBlowsChart.data.fullLabels = fullLabels;
            damageFinalBlowsChart.data.datasets[0].data = damageData;
            damageFinalBlowsChart.data.datasets[1].data = finalBlowsData;
            damageFinalBlowsChart.update();
        } else {
            // Create new chart
            damageFinalBlowsChart = new Chart(ctxElem.getContext('2d'), {
                type: 'bar',
                data: {
                    labels: truncatedLabels,
                    fullLabels: fullLabels,
                    datasets: [
                        {
                            label: 'Damage Done',
                            data: damageData,
                            backgroundColor: 'rgba(255, 77, 77, 0.7)',
                        },
                        {
                            label: 'Final Blows',
                            data: finalBlowsData,
                            backgroundColor: 'rgba(54, 162, 235, 0.7)',
                        },
                    ],
                },
                options: {
                    indexAxis: 'y',
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            display: true,
                        },
                        tooltip: {
                            callbacks: {
                                title: function(context) {
                                    const index = context[0].dataIndex;
                                    return context[0].chart.data.fullLabels[index];
                                },
                            },
                        },
                    },
                    scales: {
                        x: {
                            ticks: { color: '#ffffff' },
                            grid: { color: '#444' },
                            beginAtZero: true,
                        },
                        y: {
                            ticks: {
                                color: '#ffffff',
                                autoSkip: false,
                            },
                            grid: { display: false },
                        },
                    },
                },
            });
        }
    }

    // 2. Our Losses Combined Chart
    function updateOurLossesCombinedChart() {
        let data;
        if (currentTimeFrame === 'mtd') {
            data = window.mtdOurLossesValueData;
        } else if (currentTimeFrame === 'ytd') {
            data = window.ytdOurLossesValueData;
        } else if (currentTimeFrame === 'lastMonth') {
            data = window.lastMOurLossesValueData;
        }

        if (!data || !Array.isArray(data) || data.length === 0) {
            console.warn('Our Losses Combined Data is unavailable or empty for the selected time frame.');
            return;
        }

        // Filter out invalid items
        data = data.filter(item => item && item.CharacterName);

        const labels = data.map(item => item.CharacterName);
        const lossesValueData = data.map(item => item.TotalValue || 0);
        const lossesCountData = data.map(item => item.LossesCount || 0);
        const shipCountData = data.map(item => item.ShipCount || 0);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));


        const ctxElem = document.getElementById('ourLossesCombinedChart');
        if (!ctxElem) return;

        if (ourLossesCombinedChart) {
            ourLossesCombinedChart.data.labels = truncatedLabels;
            ourLossesCombinedChart.data.fullLabels = fullLabels;
            ourLossesCombinedChart.data.datasets[0].data = lossesValueData;
            ourLossesCombinedChart.data.datasets[1].data = lossesCountData;
            ourLossesCombinedChart.data.datasets[2].data = shipCountData;
            ourLossesCombinedChart.update();
        } else {
            ourLossesCombinedChart = new Chart(ctxElem.getContext('2d'), {
                type: 'bar',
                data: {
                    labels: truncatedLabels,
                    fullLabels: fullLabels,
                    datasets: [
                        {
                            label: 'Losses Value (ISK)',
                            data: lossesValueData,
                            backgroundColor: 'rgba(255, 77, 77, 0.7)',
                            yAxisID: 'y1',
                        },
                        {
                            label: 'Losses Count',
                            data: lossesCountData,
                            backgroundColor: 'rgba(54, 162, 235, 0.7)',
                            yAxisID: 'y',
                        },
                        {
                            label: 'Ship Types Lost',
                            data: shipCountData,
                            backgroundColor: 'rgba(75, 192, 192, 0.7)',
                            yAxisID: 'y',
                        },
                    ],
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            display: true,
                        },
                        tooltip: {
                            mode: 'index',
                            intersect: false,
                            callbacks: {
                                title: function(context) {
                                    const index = context[0].dataIndex;
                                    return context[0].chart.data.fullLabels[index];
                                },
                            },
                        },
                    },
                    scales: {
                        y: {
                            type: 'linear',
                            position: 'left',
                            ticks: { color: '#ffffff' },
                            grid: { color: '#444' },
                            beginAtZero: true,
                        },
                        y1: {
                            type: 'linear',
                            position: 'right',
                            ticks: { color: '#ffffff' },
                            grid: { drawOnChartArea: false },
                            beginAtZero: true,
                        },
                        x: {
                            ticks: {
                                color: '#ffffff',
                                maxRotation: 45,
                                minRotation: 45,
                                autoSkip: false,
                            },
                            grid: { display: false },
                        },
                    },
                },
            });
        }
    }

    // 3. Character Performance Chart
    function updateCharacterPerformanceChart() {
        let data;
        if (currentTimeFrame === 'mtd') {
            data = window.mtdCharacterPerformanceData;
        } else if (currentTimeFrame === 'ytd') {
            data = window.ytdCharacterPerformanceData;
        } else if (currentTimeFrame === 'lastMonth') {
            data = window.lastMCharacterPerformanceData;
        }

        if (!data) return;

        const labels = data.map(item => item.CharacterName);
        const killCountData = data.map(item => item.KillCount);
        const soloKillsData = data.map(item => item.SoloKills);
        const pointsData = data.map(item => item.Points);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const ctxElem = document.getElementById('characterPerformanceChart');
        if (!ctxElem) return;

        if (characterPerformanceChart) {
            characterPerformanceChart.data.labels = truncatedLabels;
            characterPerformanceChart.data.fullLabels = fullLabels;
            characterPerformanceChart.data.datasets[0].data = killCountData;
            characterPerformanceChart.data.datasets[1].data = soloKillsData;
            characterPerformanceChart.data.datasets[2].data = pointsData;
            characterPerformanceChart.update();
        } else {
            characterPerformanceChart = new Chart(ctxElem.getContext('2d'), {
                data: {
                    labels: truncatedLabels,
                    fullLabels: fullLabels,
                    datasets: [
                        {
                            label: 'Kill Count',
                            data: killCountData,
                            backgroundColor: 'rgba(255, 77, 77, 0.7)',
                            yAxisID: 'y',
                            type: 'bar',
                        },
                        {
                            label: 'Solo Kills',
                            data: soloKillsData,
                            backgroundColor: 'rgba(54, 162, 235, 0.7)',
                            yAxisID: 'y',
                            type: 'bar',
                        },
                        {
                            label: 'Points',
                            data: pointsData,
                            borderColor: 'rgba(255, 206, 86, 1)',
                            backgroundColor: 'rgba(255, 206, 86, 0.5)',
                            yAxisID: 'y1',
                            type: 'line',
                        },
                    ],
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            display: true,
                        },
                        tooltip: {
                            mode: 'index',
                            intersect: false,
                            callbacks: {
                                title: function(context) {
                                    const index = context[0].dataIndex;
                                    return context[0].chart.data.fullLabels[index];
                                },
                            },
                        },
                    },
                    scales: {
                        y: {
                            type: 'linear',
                            position: 'left',
                            beginAtZero: true,
                            ticks: { color: '#ffffff' },
                        },
                        y1: {
                            type: 'linear',
                            position: 'right',
                            beginAtZero: true,
                            ticks: { color: '#ffffff' },
                            grid: { drawOnChartArea: false },
                        },
                        x: {
                            ticks: {
                                color: '#ffffff',
                                maxRotation: 45,
                                minRotation: 45,
                                autoSkip: false,
                            },
                            grid: { display: false },
                        },
                    },
                },
            });
        }
    }

    // 4. Our Ships Used Chart
    function updateOurShipsUsedChart() {
        let data;
        if (currentTimeFrame === 'mtd') {
            data = window.mtdOurShipsUsedData;
        } else if (currentTimeFrame === 'ytd') {
            data = window.ytdOurShipsUsedData;
        } else if (currentTimeFrame === 'lastMonth') {
            data = window.lastMOurShipsUsedData;
        }

        if (!data) return;

        const characters = data.Characters;
        const shipNames = data.ShipNames;
        const seriesData = data.SeriesData;

        const fullLabels = [...characters];
        const truncatedLabels = characters.map(label => truncateLabel(label, 10));

        const datasets = shipNames.map((shipName, index) => ({
            label: shipName,
            data: seriesData[shipName],
            backgroundColor: getColor(index),
        }));

        const ctxElem = document.getElementById('ourShipsUsedChart');
        if (!ctxElem) return;

        if (ourShipsUsedChart) {
            ourShipsUsedChart.data.labels = truncatedLabels;
            ourShipsUsedChart.data.fullLabels = fullLabels;
            ourShipsUsedChart.data.datasets = datasets;
            ourShipsUsedChart.update();
        } else {
            ourShipsUsedChart = new Chart(ctxElem.getContext('2d'), {
                type: 'bar',
                data: {
                    labels: truncatedLabels,
                    fullLabels: fullLabels,
                    datasets: datasets,
                },
                options: {
                    indexAxis: 'y',
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        tooltip: {
                            callbacks: {
                                title: function(context) {
                                    const index = context[0].dataIndex;
                                    return context[0].chart.data.fullLabels[index];
                                },
                                label: function(context) {
                                    const shipName = context.dataset.label;
                                    const value = context.parsed.x;
                                    return `${shipName}: ${value}`;
                                },
                            },
                        },
                        legend: {
                            display: true,
                        },
                    },
                    scales: {
                        x: {
                            stacked: true,
                            ticks: { color: '#ffffff' },
                            grid: { color: '#444' },
                        },
                        y: {
                            stacked: true,
                            ticks: {
                                color: '#ffffff',
                                autoSkip: false,
                            },
                            grid: { display: false },
                        },
                    },
                },
            });
        }
    }

    // 5. Victims Sunburst Chart
    function updateVictimsSunburstChart() {
        let data;
        if (currentTimeFrame === 'mtd') {
            data = window.mtdVictimsSunburstData;
        } else if (currentTimeFrame === 'ytd') {
            data = window.ytdVictimsSunburstData;
        } else if (currentTimeFrame === 'lastMonth') {
            data = window.lastMVictimsSunburstData;
        }

        if (!data) return;

        const ctxElem = document.getElementById('victimsSunburstChart');
        if (!ctxElem) return;

        if (victimsSunburstChart) {
            victimsSunburstChart.data.datasets[0].data = data;
            victimsSunburstChart.update();
        } else {
            victimsSunburstChart = new Chart(ctxElem.getContext('2d'), {
                type: 'sunburst',
                data: {
                    datasets: [{
                        data: data,
                        backgroundColor: function(context) {
                            const index = context.dataIndex;
                            return getColor(index);
                        },
                    }],
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        tooltip: {
                            callbacks: {
                                label: function(context) {
                                    const label = context.raw.name || '';
                                    const value = context.raw.value || 0;
                                    return `${label}: ${value}`;
                                },
                            },
                        },
                        legend: { display: false },
                    },
                },
            });
        }
    }

    // 6. Kill Activity Over Time Chart
    function updateKillActivityChart() {
        let data;
        if (currentTimeFrame === 'mtd') {
            data = window.mtdKillActivityData;
        } else if (currentTimeFrame === 'ytd') {
            data = window.ytdKillActivityData;
        } else if (currentTimeFrame === 'lastMonth') {
            data = window.lastMKillActivityData;
        }

        if (!data) return;

        const labels = data.map(item => new Date(item.Time).toLocaleDateString());
        const kills = data.map(item => item.Value);

        const ctxElem = document.getElementById('killActivityChart');
        if (!ctxElem) return;

        if (killActivityChart) {
            killActivityChart.data.labels = labels;
            killActivityChart.data.datasets[0].data = kills;
            killActivityChart.update();
        } else {
            killActivityChart = new Chart(ctxElem.getContext('2d'), {
                type: 'line',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Kills Over Time',
                        data: kills,
                        borderColor: 'rgba(255, 77, 77, 1)',
                        backgroundColor: 'rgba(255, 77, 77, 0.5)',
                        fill: true,
                    }],
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: { display: false },
                    },
                    scales: {
                        x: {
                            type: 'time',
                            time: {
                                unit: 'day',
                            },
                            ticks: { color: '#ffffff' },
                            grid: { color: '#444' },
                        },
                        y: {
                            beginAtZero: true,
                            ticks: { color: '#ffffff' },
                            grid: { color: '#444' },
                        },
                    },
                },
            });
        }
    }

    // 7. Kill Heatmap Chart
    function updateKillHeatmapChart() {
        let data;
        if (currentTimeFrame === 'mtd') {
            data = window.mtdKillHeatmapData;
        } else if (currentTimeFrame === 'ytd') {
            data = window.ytdKillHeatmapData;
        } else if (currentTimeFrame === 'lastMonth') {
            data = window.lastMKillHeatmapData;
        }

        if (!data) return;

        const maxKills = Math.max(...data.flat());
        const heatmapData = [];

        for (let day = 0; day < 7; day++) {
            for (let hour = 0; hour < 24; hour++) {
                const kills = data[day][hour];
                heatmapData.push({
                    x: hour,
                    y: day,
                    v: kills,
                });
            }
        }

        const ctxElem = document.getElementById('killHeatmapChart');
        if (!ctxElem) return;

        if (killHeatmapChart) {
            killHeatmapChart.data.datasets[0].data = heatmapData;
            killHeatmapChart.update();
        } else {
            killHeatmapChart = new Chart(ctxElem.getContext('2d'), {
                type: 'matrix',
                data: {
                    datasets: [{
                        label: 'Kill Heatmap',
                        data: heatmapData,
                        backgroundColor: function(ctx) {
                            const value = ctx.dataset.data[ctx.dataIndex].v;
                            const alpha = value / maxKills;
                            return `rgba(255, 77, 77, ${alpha})`;
                        },
                        width: ({ chart }) => (chart.chartArea || {}).width / 24 - 1,
                        height: ({ chart }) => (chart.chartArea || {}).height / 7 - 1,
                    }],
                },
                options: {
                    scales: {
                        x: {
                            type: 'category',
                            labels: [...Array(24).keys()],
                            ticks: { color: '#ffffff' },
                            grid: { display: false },
                        },
                        y: {
                            type: 'category',
                            labels: ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'],
                            ticks: { color: '#ffffff' },
                            grid: { display: false },
                        },
                    },
                    plugins: {
                        tooltip: {
                            callbacks: {
                                label: function(context) {
                                    const x = context.raw.x;
                                    const y = context.raw.y;
                                    const value = context.raw.v;
                                    const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
                                    return `Day: ${days[y]}, Hour: ${x}, Kills: ${value}`;
                                },
                            },
                        },
                        legend: { display: false },
                    },
                },
            });
        }
    }

    // 8. Kill-to-Loss Ratio Chart
    function updateKillLossRatioChart() {
        let data;
        if (currentTimeFrame === 'mtd') {
            data = window.mtdKillLossRatioData;
        } else if (currentTimeFrame === 'ytd') {
            data = window.ytdKillLossRatioData;
        } else if (currentTimeFrame === 'lastMonth') {
            data = window.lastMKillLossRatioData;
        }

        if (!data) return;

        const labels = data.map(item => item.CharacterName);
        const ratios = data.map(item => item.Ratio);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const ctxElem = document.getElementById('killLossRatioChart');
        if (!ctxElem) return;

        if (killLossRatioChart) {
            killLossRatioChart.data.labels = truncatedLabels;
            killLossRatioChart.data.fullLabels = fullLabels;
            killLossRatioChart.data.datasets[0].data = ratios;
            killLossRatioChart.update();
        } else {
            killLossRatioChart = new Chart(ctxElem.getContext('2d'), {
                type: 'bar',
                data: {
                    labels: truncatedLabels,
                    fullLabels: fullLabels,
                    datasets: [{
                        label: 'Kill-to-Loss Ratio',
                        data: ratios,
                        backgroundColor: 'rgba(75, 192, 192, 0.7)',
                    }],
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: { display: false },
                        tooltip: {
                            callbacks: {
                                title: function(context) {
                                    const index = context[0].dataIndex;
                                    return context[0].chart.data.fullLabels[index];
                                },
                                label: function(context) {
                                    const kills = data[context.dataIndex].Kills;
                                    const losses = data[context.dataIndex].Losses;
                                    const ratio = context.parsed.y.toFixed(2);
                                    return `Kills: ${kills}, Losses: ${losses}, Ratio: ${ratio}`;
                                },
                            },
                        },
                    },
                    scales: {
                        x: {
                            ticks: {
                                color: '#ffffff',
                                maxRotation: 45,
                                minRotation: 45,
                                autoSkip: false,
                            },
                            grid: { display: false },
                        },
                        y: {
                            beginAtZero: true,
                            ticks: { color: '#ffffff' },
                            grid: { color: '#444' },
                        },
                    },
                },
            });
        }
    }

    // 9. Top Ships Killed Chart
    function updateTopShipsKilledChart() {
        let data;
        if (currentTimeFrame === 'mtd') {
            data = window.mtdTopShipsKilledData;
        } else if (currentTimeFrame === 'ytd') {
            data = window.ytdTopShipsKilledData;
        } else if (currentTimeFrame === 'lastMonth') {
            data = window.lastMTopShipsKilledData;
        }

        if (!data) return;

        const labels = data.map(item => item.ShipName);
        const killCounts = data.map(item => item.KillCount);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 15));

        const ctxElem = document.getElementById('topShipsKilledChart');
        if (!ctxElem) return;

        if (topShipsKilledChart) {
            topShipsKilledChart.data.labels = truncatedLabels;
            topShipsKilledChart.data.fullLabels = fullLabels;
            topShipsKilledChart.data.datasets[0].data = killCounts;
            topShipsKilledChart.update();
        } else {
            topShipsKilledChart = new Chart(ctxElem.getContext('2d'), {
                type: 'bar',
                data: {
                    labels: truncatedLabels,
                    fullLabels: fullLabels,
                    datasets: [{
                        label: 'Ships Killed',
                        data: killCounts,
                        backgroundColor: 'rgba(255, 77, 77, 0.7)',
                    }],
                },
                options: {
                    indexAxis: 'y',
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: { display: false },
                        tooltip: {
                            callbacks: {
                                title: function(context) {
                                    const index = context[0].dataIndex;
                                    return context[0].chart.data.fullLabels[index];
                                },
                            },
                        },
                    },
                    scales: {
                        x: {
                            ticks: { color: '#ffffff' },
                            grid: { color: '#444' },
                            beginAtZero: true,
                        },
                        y: {
                            ticks: {
                                color: '#ffffff',
                                autoSkip: false,
                            },
                            grid: { display: false },
                        },
                    },
                },
            });
        }
    }

    // 10. Value Over Time Chart
    function updateValueOverTimeChart() {
        let data;
        if (currentTimeFrame === 'mtd') {
            data = window.mtdValueOverTimeData;
        } else if (currentTimeFrame === 'ytd') {
            data = window.ytdValueOverTimeData;
        } else if (currentTimeFrame === 'lastMonth') {
            data = window.lastMValueOverTimeData;
        }

        if (!data) return;

        const labels = data.map(item => new Date(item.Time).toLocaleDateString());
        const values = data.map(item => item.Value);

        const ctxElem = document.getElementById('valueOverTimeChart');
        if (!ctxElem) return;

        if (valueOverTimeChart) {
            valueOverTimeChart.data.labels = labels;
            valueOverTimeChart.data.datasets[0].data = values;
            valueOverTimeChart.update();
        } else {
            valueOverTimeChart = new Chart(ctxElem.getContext('2d'), {
                type: 'line',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'ISK Value Destroyed Over Time',
                        data: values,
                        borderColor: 'rgba(54, 162, 235, 1)',
                        backgroundColor: 'rgba(54, 162, 235, 0.5)',
                        fill: true,
                    }],
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: { display: false },
                    },
                    scales: {
                        x: {
                            type: 'time',
                            time: {
                                unit: 'day',
                            },
                            ticks: { color: '#ffffff' },
                            grid: { color: '#444' },
                        },
                        y: {
                            beginAtZero: true,
                            ticks: { color: '#ffffff' },
                            grid: { color: '#444' },
                        },
                    },
                },
            });
        }
    }

    // Initial chart rendering
    updateAllCharts();
});
