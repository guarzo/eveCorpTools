// static/js/tps.js

// Import utility functions
import { truncateLabel, getColor, getCommonOptions } from './utils.js';

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

// Array of all chart configurations
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
     * @returns {Chart} - The initialized Chart.js instance.
     */
    function createChart(config, ctxElem, data) {
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
     */
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

    /**
     * Updates or creates a chart based on the current time frame.
     * @param {Object} config - The chart configuration object.
     */
    function updateChart(config) {
        const dataKey = config.dataKeys[currentTimeFrame];
        let data = window[dataKey];

        if (!data || (Array.isArray(data) && data.length === 0)) {
            console.warn(`Data unavailable for chart ${config.id} in ${currentTimeFrame}.`);
            if (config.instance) {
                config.instance.destroy();
                config.instance = null;
            }
            // Optionally, display a placeholder or message
            const ctxElem = document.getElementById(config.id);
            if (ctxElem) {
                ctxElem.parentElement.innerHTML = `<p class="text-center">No data available for this chart.</p>`;
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
