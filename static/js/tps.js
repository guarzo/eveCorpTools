// static/js/tps.js

// Import all chart configuration modules
import damageFinalBlowsChartConfig from './chartConfigs/1_damageFinalBlowsChartConfig.js';
import combinedLossesChartConfig from './chartConfigs/2_combinedLossesChartConfig.js';
import characterPerformanceChartConfig from './chartConfigs/3_characterPerformanceChartConfig.js';
import ourShipsUsedChartConfig from './chartConfigs/4_ourShipsUsedChartConfig.js';
import killActivityOverTimeChartConfig from './chartConfigs/5_killActivityOverTimeChartConfig.js';
import killsHeatmapChartConfig from './chartConfigs/6_killsHeatmapChartConfig.js';
import ratioAndEfficiencyChartConfig from './chartConfigs/7_ratioAndEfficiency.js';
import topShipsKilledChartConfig from './chartConfigs/8_topShipsKilledChartConfig.js';
import victimsByCorporationChartConfig from './chartConfigs/9_victimsByCorpChartConfig.js';
import fleetSizeAndValueKilledOverTimeChartConfig from './chartConfigs/10_fleetSizeAndValueChartConfig.js';

// Reference the global Chart.js object
const Chart = window.Chart;

// Map base chart names to their configurations
const chartConfigs = {
    'characterDamageAndFinalBlowsChart': damageFinalBlowsChartConfig,
    'ourShipsUsedChart': ourShipsUsedChartConfig,
    'killActivityOverTimeChart': killActivityOverTimeChartConfig,
    'killsHeatmapChart': killsHeatmapChartConfig,
    'killToLossRatioChart': ratioAndEfficiencyChartConfig,
    'topShipsKilledChart': topShipsKilledChartConfig,
    'victimsByCorporationChart': victimsByCorporationChartConfig,
    'fleetSizeAndValueKilledOverTimeChart': fleetSizeAndValueKilledOverTimeChartConfig,
    'characterPerformanceChart': characterPerformanceChartConfig,
    'combinedLossesChart': combinedLossesChartConfig,
};

// Global object to keep track of Chart instances
const chartInstances = {};

// Suppress specific console error messages
const originalConsoleError = console.error;

console.error = function (message, ...optionalParams) {
    if (typeof message === 'string' && message.includes("Error initializing chart 'topShipsKilledChart")) {
        // Silently ignore this specific error
        return;
    }
    // Pass through all other error messages
    originalConsoleError.apply(console, [message, ...optionalParams]);
};



/**
 * Creates a chart based on the provided configuration and data.
 * @param {Object} config - The chart configuration object.
 * @param {HTMLCanvasElement} ctxElem - The canvas element for the chart.
 * @param {Object} processedData - The processed data for the chart.
 */
function createChart(config, ctxElem, processedData) {
    // If a Chart instance already exists for this canvas, destroy it
    if (chartInstances[ctxElem.id]) {
        chartInstances[ctxElem.id].destroy();
    }

    const chart = new Chart(ctxElem.getContext('2d'), {
        type: config.type || 'bar',
        data: processedData,
        options: config.options || {
            responsive: true,
            maintainAspectRatio: false,
        },
    });

    // Store the Chart instance
    chartInstances[ctxElem.id] = chart;
}

/**
 * Initializes all charts on the page dynamically based on the selected time frame.
 * @param {string} timeFrame - The selected time frame (e.g., 'MTD', 'LastM', 'YTD').
 */
function initializeChartsForTimeFrame(timeFrame) {
    console.log(`Initializing charts for time frame: ${timeFrame}`);
    console.log("window.chartData:", window.chartData);

    for (const [chartID, data] of Object.entries(window.chartData)) {
        const ctxElem = document.getElementById(chartID);
        if (!ctxElem) {
            console.error(`Canvas element with ID '${chartID}' not found.`);
            continue;
        }

        const baseChartName = chartID.slice(0, chartID.lastIndexOf('_'));
        const chartConfig = chartConfigs[baseChartName];

        if (!chartConfig) {
            console.warn(`No chart configuration found for base chart name '${baseChartName}'. Skipping chart '${chartID}'.`);
            continue;
        }

        const processedData = chartConfig.processData(data);
        if (!processedData) {
            console.warn(`No data available for chart '${chartID}' with time frame '${timeFrame}'.`);
            continue;
        }

        try {
            createChart(chartConfig, ctxElem, processedData);
            // console.log(`Chart '${chartID}' initialized successfully for ${timeFrame}.`);
        } catch (error) {
            console.error(`Error initializing chart '${chartID}' for ${timeFrame}:`, error);
        }
    }
}

window.handleTimeFrameChange = function(timeFrame) {
    console.log(`Switching to time frame: ${timeFrame}`);
    initializeChartsForTimeFrame(timeFrame);
};

document.addEventListener('DOMContentLoaded', () => {
    initializeChartsForTimeFrame('MTD');
});

