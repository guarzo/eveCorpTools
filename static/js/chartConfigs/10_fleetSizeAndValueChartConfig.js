// static/js/chartConfigs/10_fleetSizeAndValueKilledOverTimeChartConfig.js

import { getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Fleet Size and Average Value Killed Over Time Chart
 */
const fleetSizeAndValueKilledOverTimeChartConfig = {
    type: 'line', // Using 'line' chart type
    options: getCommonOptions('Fleet Size and Average Value Killed Over Time', {
        // ... your specific options here ...
        scales: {
            x: {
                title: {
                    display: true,
                    text: 'Time',
                },
                ticks: {
                    color: '#ffffff',
                },
                grid: { display: false },
            },
            y: {
                title: {
                    display: true,
                    text: 'Value Killed',
                },
                ticks: {
                    color: '#ffffff',
                },
                grid: { display: false },
            },
        },
        plugins: {
            legend: {
                display: true,
                position: 'top',
                labels: {
                    color: '#ffffff',
                },
            },
            tooltip: {
                callbacks: {
                    label: function(context) {
                        const label = context.dataset.label || '';
                        const value = context.parsed.y || 0;
                        return `${label}: ${value}`;
                    },
                },
            },
        },
    }),
    processData: function(data) {
        const chartName = 'Fleet Size and Average Value Killed Over Time';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty data to trigger the noDataPlugin
            return { labels: [], datasets: [], noDataMessage: 'No data available for this chart.' };
        }

        console.log('Incoming data for Fleet Size and Average Value Killed Over Time:', data); // Debugging log

        // Inspect each data item
        data.forEach((item, index) => {
            console.log(`Item ${index}:`, item);
        });

        // Update field names based on actual data structure
        const labels = data.map(item => item.timePeriod || item.time || 'Unknown');
        const fleetSizes = data.map(item => item.fleetSize || 0);
        const averageValues = data.map(item => item.averageValueKilled || 0);

        // Check for 'Unknown' labels
        const allUnknown = labels.every(label => label === 'Unknown');
        if (allUnknown) {
            console.warn(`All labels for ${chartName} are 'Unknown'. Check data source.`);
        }

        return {
            labels: labels,
            datasets: [
                {
                    label: 'Fleet Size Killed',
                    data: fleetSizes,
                    borderColor: 'rgba(75, 192, 192, 1)',
                    backgroundColor: 'rgba(75, 192, 192, 0.2)',
                    fill: false,
                    tension: 0.1, // Smoothness of the line
                },
                {
                    label: 'Average Value Killed',
                    data: averageValues,
                    borderColor: 'rgba(153, 102, 255, 1)',
                    backgroundColor: 'rgba(153, 102, 255, 0.2)',
                    fill: false,
                    tension: 0.1,
                },
            ],
        };
    },
};

export default fleetSizeAndValueKilledOverTimeChartConfig;
