// tps.js

// Import chart configurations
import damageFinalBlowsChartConfig from './chartConfigs/damageFinalBlowsChartConfig.js';
import ourLossesCombinedChartConfig from './chartConfigs/ourLossesCombinedChartConfig.js';
import characterPerformanceChartConfig from './chartConfigs/characterPerformanceChartConfig.js';
import ourShipsUsedChartConfig from './chartConfigs/ourShipsUsedChartConfig.js';
import victimsSunburstChartConfig from './chartConfigs/victimsSunburstChartConfig.js';
import killActivityChartConfig from './chartConfigs/killActivityChartConfig.js';
import killHeatmapChartConfig from './chartConfigs/killHeatmapChartConfig.js';
import killLossRatioChartConfig from './chartConfigs/killLossRatioChartConfig.js';
import topShipsKilledChartConfig from './chartConfigs/topShipsKilledChartConfig.js';
import valueOverTimeChartConfig from './chartConfigs/valueOverTimeChartConfig.js';


let currentTimeFrame = 'mtd';

// Set initial active button
const btnMtd = document.getElementById('btn-mtd');
if (btnMtd) {
    btnMtd.classList.add('active');
} else {
    console.error('Element with id "btn-mtd" not found.');
}

// Function to set the time frame
window.setTimeFrame = function (timeFrame) {
    currentTimeFrame = timeFrame;
    updateAllCharts();

    // Update button styles
    ['mtd', 'lastMonth', 'ytd'].forEach(tf => {
        const btn = document.getElementById(`btn-${tf}`);
        if (btn) {
            btn.classList.toggle('active', tf === timeFrame);
        }
    });
};

// List of chart configurations
const chartConfigs = [
    damageFinalBlowsChartConfig,
    ourLossesCombinedChartConfig,
    characterPerformanceChartConfig,
    ourShipsUsedChartConfig,
    victimsSunburstChartConfig,
    killActivityChartConfig,
    killHeatmapChartConfig,
    killLossRatioChartConfig,
    topShipsKilledChartConfig,
    valueOverTimeChartConfig,
];

// Function to create charts
function createChart(config, ctxElem, data) {
    const { labels, datasets, fullLabels, originalData } = data;

    return new Chart(ctxElem.getContext('2d'), {
        type: config.type,
        data: {
            labels: labels,
            datasets: datasets,
            fullLabels: fullLabels,
            originalData: originalData,
        },
        options: config.options,
    });
}

// Function to update existing charts
function updateChartInstance(config, data) {
    const { labels, datasets, fullLabels } = data;
    const chart = config.instance;

    chart.data.labels = labels;
    chart.data.datasets = datasets;
    if (fullLabels) {
        chart.data.fullLabels = fullLabels;
    }
    chart.update();
}

// Function to update or create charts
function updateChart(config) {
    const dataKey = config.dataKeys[currentTimeFrame];
    let data = window[dataKey];

    if (!data || (Array.isArray(data) && data.length === 0)) {
        console.warn(`Data unavailable for chart ${config.id} in ${currentTimeFrame}.`);
        if (config.instance) {
            config.instance.destroy();
            config.instance = null;
        }
        return;
    }

    if (!Array.isArray(data)) {
        data = [data];
    }

    data = data.filter(item => item);

    const processedData = config.processData(data);
    if (!processedData) {
        console.warn(`Processed data is invalid for chart ${config.id}.`);
        return;
    }

    const ctxElem = document.getElementById(config.id);
    if (!ctxElem) return;

    if (config.instance) {
        updateChartInstance(config, processedData);
    } else {
        config.instance = createChart(config, ctxElem, processedData);
    }
}

// Function to update all charts
function updateAllCharts() {
    chartConfigs.forEach(config => {
        updateChart(config);
    });
}

// Initial chart rendering
updateAllCharts();
