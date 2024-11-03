// static/js/tps.js

// Import all chart configuration modules
import damageFinalBlowsChartConfig from './chartConfigs/1_damageFinalBlowsChartConfig.js';
import combinedLossesChartConfig from './chartConfigs/2_combinedLossesChartConfig.js';
import characterPerformanceChartConfig from './chartConfigs/3_characterPerformanceChartConfig.js';
import ourShipsUsedChartConfig from './chartConfigs/4_ourShipsUsedChartConfig.js';
import killActivityOverTimeChartConfig from './chartConfigs/5_killActivityOverTimeChartConfig.js';
import killsHeatmapChartConfig from './chartConfigs/6_killsHeatmapChartConfig.js';
import killToLossRatioChartConfig from './chartConfigs/7_killLossRatioChartConfig.js';
import topShipsKilledChartConfig from './chartConfigs/8_topShipsKilledChartConfig.js';
import victimsByCorporationChartConfig from './chartConfigs/9_victimsByCorpChartConfig.js';
import fleetSizeAndValueKilledOverTimeChartConfig from './chartConfigs/10_fleetSizeAndValueChartConfig.js';

// Reference the global Chart.js object
const Chart = window.Chart;

// Map base chart names to their configurations
const chartConfigs = {
    'characterDamageAndFinalBlowsChart': damageFinalBlowsChartConfig,
    'combinedLossesChart': combinedLossesChartConfig,
    'characterPerformanceChart': characterPerformanceChartConfig,
    'ourShipsUsedChart': ourShipsUsedChartConfig,
    'killActivityOverTimeChart': killActivityOverTimeChartConfig,
    'killsHeatmapChart': killsHeatmapChartConfig,
    'killToLossRatioChart': killToLossRatioChartConfig,
    'topShipsKilledChart': topShipsKilledChartConfig,
    'victimsByCorporationChart': victimsByCorporationChartConfig,
    'fleetSizeAndValueKilledOverTimeChart': fleetSizeAndValueKilledOverTimeChartConfig,
};

// Global object to keep track of Chart instances
const chartInstances = {};

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

    let chart;
    switch (config.type.toLowerCase()) {
        case 'wordcloud':
            // Handle Word Cloud Charts
            chart = new Chart(ctxElem.getContext('2d'), {
                type: 'wordCloud',
                data: {
                    labels: processedData.labels, // Not used in Word Cloud
                    datasets: [{
                        data: processedData.datasets[0].data, // Array of { text, value }
                        backgroundColor: processedData.datasets[0].backgroundColor,
                        rotation: processedData.datasets[0].rotation,
                        weightFactor: processedData.datasets[0].weightFactor,
                    }]
                },
                options: config.options,
            });
            break;

        case 'matrix':
            // Handle Matrix Charts
            chart = new Chart(ctxElem.getContext('2d'), {
                type: 'matrix',
                data: {
                    datasets: [{
                        label: 'Matrix Data',
                        data: processedData.datasets[0].data, // Array of { x, y, v }
                        backgroundColor: processedData.datasets[0].backgroundColor,
                        borderWidth: 1,
                        borderColor: '#ffffff',
                        borderSkipped: 'bottom',
                        width: function(context) {
                            const a = context.chart.chartArea;
                            return (a && a.right && a.left) ? (a.right - a.left) / 24 : 20; // 24 hours
                        },
                        height: function(context) {
                            const a = context.chart.chartArea;
                            return (a && a.bottom && a.top) ? (a.bottom - a.top) / 7 : 20; // 7 days
                        },
                    }]
                },
                options: config.options,
            });
            break;

        default:
            // Handle standard chart types (bar, line, etc.)
            chart = new Chart(ctxElem.getContext('2d'), {
                type: config.type || 'bar',
                data: processedData,
                options: config.options || {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: [], // Add any global plugins if necessary
                },
            });
            break;
    }

    // Store the Chart instance
    chartInstances[ctxElem.id] = chart;
}

/**
 * Initializes all charts on the page dynamically based on window.chartData and chartConfigs.
 */
function init() {
    /**
     * Initializes all charts based on window.chartData and chartConfigs.
     */
    function initializeCharts() {
        console.log('Initializing Charts...');
        for (const [chartID, data] of Object.entries(window.chartData)) {
            const ctxElem = document.getElementById(chartID);
            if (!ctxElem) {
                console.error(`Canvas element with ID '${chartID}' not found.`);
                continue;
            }

            // Extract the base chart name by removing the timeframe suffix
            // Example: 'damageFinalBlowsChart_MTD' => 'damageFinalBlowsChart'
            const lastUnderscore = chartID.lastIndexOf('_');
            if (lastUnderscore === -1) {
                console.warn(`Chart ID '${chartID}' does not contain a timeframe suffix. Skipping.`);
                continue;
            }
            const baseChartName = chartID.substring(0, lastUnderscore);
            const chartConfig = chartConfigs[baseChartName];

            if (!chartConfig) {
                console.warn(`No chart configuration found for base chart name '${baseChartName}'. Skipping chart '${chartID}'.`);
                continue;
            }

            // Process the data using the chart's processData function
            const processedData = chartConfig.processData(data);

            if (!processedData || (processedData.labels.length === 0 && processedData.datasets.length === 0)) {
                console.warn(`Processed data for chart '${chartID}' is empty. Skipping initialization.`);
                continue;
            }

            // Create the chart
            try {
                createChart(chartConfig, ctxElem, processedData);
                console.log(`Chart '${chartID}' initialized successfully.`);
            } catch (error) {
                console.error(`Error initializing chart '${chartID}':`, error);
            }
        }
    }

    /**
     * Handles time frame changes by re-initializing charts.
     * @param {string} timeFrame - The selected time frame (e.g., 'MTD', 'LastM', 'YTD').
     */
    function handleTimeFrameChange(timeFrame) {
        // Optional: Implement any additional logic needed when the timeframe changes
        // Currently, charts are re-initialized with new data
        initializeCharts();
    }

    /**
     * Sets up event listeners for Bootstrap tabs to handle timeframe changes.
     */
    function setupTabListeners() {
        const chartTab = document.getElementById('chartTab');
        if (chartTab) {
            chartTab.addEventListener('shown.bs.tab', function (event) {
                const activatedTabId = event.target.id; // e.g., 'mtd-tab'
                if (activatedTabId === 'mtd-tab') {
                    handleTimeFrameChange('MTD');
                } else if (activatedTabId === 'lastm-tab') {
                    handleTimeFrameChange('LastM');
                } else if (activatedTabId === 'ytd-tab') {
                    handleTimeFrameChange('YTD');
                }
            });
        }
    }

    // Initialize charts and set up event listeners
    initializeCharts();
    setupTabListeners();
}

// Check if the DOM is already loaded
if (document.readyState === 'loading') {
    // DOM not ready, wait for it
    document.addEventListener('DOMContentLoaded', init);
} else {
    // DOM is ready
    init();
}
