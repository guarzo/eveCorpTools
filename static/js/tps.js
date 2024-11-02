// static/js/tps.js

// Import utility functions
import { truncateLabel, getColor, getCommonOptions, noDataPlugin } from './utils.js';
import { initTopShipsKilledWordCloud } from './chartConfigs/8_topShipsKilledWordCloudConfig.js';
// Reference the global Chart.js object
const Chart = window.Chart;

// Register necessary plugins
Chart.register(noDataPlugin);
// If you are using other plugins like datalabels, register them as well
// import ChartDataLabels from 'chartjs-plugin-datalabels';
// Chart.register(ChartDataLabels);

// Import chart configurations
import damageFinalBlowsChartConfig from './chartConfigs/1_damageFinalBlowsChartConfig.js';
import ourLossesCombinedChartConfig from './chartConfigs/2_ourLossesCombinedChartConfig.js';
import characterPerformanceChartConfig from './chartConfigs/3_characterPerformanceChartConfig.js';
import ourShipsUsedChartConfig from './chartConfigs/4_ourShipsUsedChartConfig.js';
import killActivityChartConfig from './chartConfigs/5_killActivityChartConfig.js';
import killHeatmapChartConfig from './chartConfigs/6_killsHeatmapChartConfig.js';
import killLossRatioChartConfig from './chartConfigs/7_killLossRatioChartConfig.js';
import victimsByCorpChartConfig from "./chartConfigs/9_victimsByCorpChartConfig.js";
import averageFleetSizeChartConfig from "./chartConfigs/10_fleetSizeChartMap.js";

// Array of all chart configurations
const chartConfigs = [
    damageFinalBlowsChartConfig,
    ourLossesCombinedChartConfig,
    characterPerformanceChartConfig,
    ourShipsUsedChartConfig,
    killActivityChartConfig,
    killHeatmapChartConfig,
    killLossRatioChartConfig,
    victimsByCorpChartConfig,
    averageFleetSizeChartConfig,
];

/**
 * Initializes all charts on the page.
 */
function init() {
    let currentTimeFrame = 'mtd';

    // Determine the initially active tab
    const activeTab = document.querySelector('.nav-tabs .nav-link.active');
    if (activeTab) {
        const activatedTabId = activeTab.id;
        if (activatedTabId === 'mtd-tab') {
            currentTimeFrame = 'mtd';
        } else if (activatedTabId === 'lm-tab') {
            currentTimeFrame = 'lastMonth';
        } else if (activatedTabId === 'ytd-tab') {
            currentTimeFrame = 'ytd';
        }
    }

    // Function to set the time frame
    function setTimeFrame(timeFrame) {
        currentTimeFrame = timeFrame;
        updateAllCharts();
    }

    // Listen for Bootstrap tab events
    const chartTab = document.getElementById('chartTab');
    if (chartTab) {
        chartTab.addEventListener('shown.bs.tab', function (event) {
            const activatedTabId = event.target.id; // activated tab id
            if (activatedTabId === 'mtd-tab') {
                setTimeFrame('mtd');
            } else if (activatedTabId === 'lm-tab') {
                setTimeFrame('lastMonth');
            } else if (activatedTabId === 'ytd-tab') {
                setTimeFrame('ytd');
            }
        });
    }

    /**
     * Creates a chart based on the provided configuration and data.
     * @param {Object} config - The chart configuration object.
     * @param {HTMLCanvasElement} ctxElem - The canvas element for the chart.
     * @param {Object} data - The processed data for the chart.
     * @param {string} timeFrame - The current time frame (mtd, ytd, lastMonth).
     * @returns {Chart} - The initialized Chart.js instance.
     */
    function createChart(config, ctxElem, data, timeFrame) {
        const { labels, datasets } = data;

        // Add backgroundColor and borderColor only if not a Word Cloud
        if (config.type !== 'wordCloud') {
            datasets.forEach(dataset => {
                dataset.backgroundColor = getColor(dataset.label);
                dataset.borderColor = getColor(dataset.label);
            });
        }

        const chartInstance = new Chart(ctxElem.getContext('2d'), {
            type: config.type,
            data: {
                labels: labels,
                datasets: datasets,
            },
            options: config.options,
        });

        console.log('Chart Instance:', chartInstance);
        return chartInstance;
    }

    /**
     * Updates an existing chart instance with new data.
     * @param {Object} config - The chart configuration object.
     * @param {Object} data - The new data for the chart.
     * @param {string} timeFrame - The current time frame (mtd, ytd, lastMonth).
     */
    function updateChartInstance(config, data, timeFrame) {
        const { datasets } = data;
        const chart = config.instance[timeFrame];

        if (chart) {
            chart.data.datasets = datasets;
            chart.update();
        }
    }

    /**
     * Updates or creates a chart based on the current time frame.
     * @param {Object} config - The chart configuration object.
     */
    function updateChart(config) {
        const dataKey = config.dataKeys[currentTimeFrame];
        let data = window[dataKey.dataVar];

        console.log(`Updating chart ${config.id} for ${currentTimeFrame}:`, data);

        if (
            !data ||
            (config.dataType === 'array' && data.length === 0) ||
            (config.dataType === 'object' && Object.keys(data).length === 0)
        ) {
            console.warn(`Data unavailable for chart ${config.id} in ${currentTimeFrame}.`);
            if (config.instance && config.instance[currentTimeFrame]) {
                config.instance[currentTimeFrame].destroy();
                config.instance[currentTimeFrame] = null;
            }
            // Optionally, display a placeholder or message
            const ctxElem = document.getElementById(dataKey.canvasId);
            if (ctxElem) {
                ctxElem.parentElement.innerHTML = `<p class="text-center">No data available for this chart.</p>`;
            }
            return;
        }

        const processedData = config.processData(data);
        if (!processedData) {
            console.warn(`Processed data is invalid for chart ${config.id}.`);
            return;
        }

        const ctxElem = document.getElementById(dataKey.canvasId);
        if (!ctxElem) {
            console.error(`Canvas element with ID '${dataKey.canvasId}' not found.`);
            return;
        }

        if (config.instance && config.instance[currentTimeFrame]) {
            // Update existing chart instance
            updateChartInstance(config, processedData, currentTimeFrame);
        } else {
            // Create a new chart instance
            config.instance = config.instance || {};
            // Clone the options to avoid mutating the original configuration
            const chartOptions = JSON.parse(JSON.stringify(config.options));

            // If there's a custom noDataMessage, set it in the options
            if (processedData.noDataMessage) {
                if (!chartOptions.plugins.noData) {
                    chartOptions.plugins.noData = {};
                }
                chartOptions.plugins.noData.message = processedData.noDataMessage;
            }

            config.instance[currentTimeFrame] = new Chart(ctxElem.getContext('2d'), {
                type: config.type,
                data: {
                    labels: processedData.labels,
                    datasets: processedData.datasets,
                },
                options: chartOptions,
            });

            console.log(`Chart Instance for ${config.id} (${currentTimeFrame}):`, config.instance[currentTimeFrame]);
        }
    }

    function createWordCloudChart(config, ctxElem, data, timeFrame) {
        const { datasets } = data;

        const chartInstance = new Chart(ctxElem.getContext('2d'), {
            type: config.type,
            data: {
                datasets: datasets,
            },
            options: config.options,
        });

        console.log('Word Cloud Chart Instance:', chartInstance);
        return chartInstance;
    }


    /**
     * Updates all charts on the page.
     */
    function updateAllCharts() {
        chartConfigs.forEach(config => {
            updateChart(config);
        });
    }

    // Initial chart rendering
    updateAllCharts();

    // Initialize the Word Cloud chart separately
    initTopShipsKilledWordCloud();
}

// Check if the DOM is already loaded
if (document.readyState === 'loading') {
    // DOM not ready, wait for it
    document.addEventListener('DOMContentLoaded', init);
} else {
    // DOM is ready
    init();
}
