// static/js/tps.js

// Import utility functions
import { truncateLabel, getColor, getCommonOptions, noDataPlugin } from './utils.js';
Chart.register(noDataPlugin);

// Import chart configurations
import damageFinalBlowsChartConfig from './chartConfigs/1_damageFinalBlowsChartConfig.js';
import ourLossesCombinedChartConfig from './chartConfigs/2_ourLossesCombinedChartConfig.js';
import characterPerformanceChartConfig from './chartConfigs/3_characterPerformanceChartConfig.js';
import ourShipsUsedChartConfig from './chartConfigs/4_ourShipsUsedChartConfig.js';
import killActivityChartConfig from './chartConfigs/killActivityChartConfig.js';
import killHeatmapChartConfig from './chartConfigs/killHeatmapChartConfig.js';
import killLossRatioChartConfig from './chartConfigs/killLossRatioChartConfig.js';
import topShipsKilledChartConfig from './chartConfigs/topShipsKilledChartConfig.js';
import valueOverTimeChartConfig from './chartConfigs/valueOverTimeChartConfig.js';

// Array of all chart configurations
const chartConfigs = [
    damageFinalBlowsChartConfig,
    ourLossesCombinedChartConfig,
    characterPerformanceChartConfig,
    ourShipsUsedChartConfig,
    killActivityChartConfig,
    killHeatmapChartConfig,
    killLossRatioChartConfig,
    topShipsKilledChartConfig,
    valueOverTimeChartConfig,
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
        const { labels, datasets, fullLabels } = data;

        return new Chart(ctxElem.getContext('2d'), {
            type: config.type,
            data: {
                labels: labels,
                datasets: datasets,
                fullLabels: fullLabels,
            },
            options: config.options,
        });
    }

    /**
     * Updates an existing chart instance with new data.
     * @param {Object} config - The chart configuration object.
     * @param {Object} data - The new data for the chart.
     * @param {string} timeFrame - The current time frame (mtd, ytd, lastMonth).
     */
    function updateChartInstance(config, data, timeFrame) {
        const { labels, datasets, fullLabels } = data;
        const chart = config.instance[timeFrame];

        if (chart) {
            chart.data.labels = labels;
            chart.data.datasets = datasets;
            if (fullLabels) {
                chart.data.fullLabels = fullLabels;
            }
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

        if (!data || (config.dataType === 'array' && data.length === 0) || (config.dataType === 'object' && Object.keys(data).length === 0)) {
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

        // Conditionally wrap data if chart expects an array
        if (config.dataType === 'array' && !Array.isArray(data)) {
            data = [data];
        }

        // Only filter if data is an array
        if (Array.isArray(data)) {
            data = data.filter(item => item);
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
            updateChartInstance(config, processedData, currentTimeFrame);
        } else {
            config.instance = config.instance || {};
            config.instance[currentTimeFrame] = createChart(config, ctxElem, processedData, currentTimeFrame);
        }
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
}

// Check if the DOM is already loaded
if (document.readyState === 'loading') {
    // DOM not ready, wait for it
    document.addEventListener('DOMContentLoaded', init);
} else {
    // DOM is ready
    init();
}
