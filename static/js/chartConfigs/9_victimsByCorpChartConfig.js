// static/js/chartConfigs/9_victimsByCorpChartConfig.js

import { getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Victims by Corporation Chart
 */
const victimsByCorporationChartConfig = {
    type: 'bar', // or the appropriate chart type
    options: getCommonOptions('Victims by Corporation', {
        // ... your specific options here ...
        scales: {
            x: {
                title: {
                    display: true,
                    text: 'Corporation',
                },
            },
            y: {
                title: {
                    display: true,
                    text: 'Number of Victims',
                },
            },
        },
    }),
    processData: function(data) {
        const chartName = 'Victims by Corporation';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty data to trigger the noDataPlugin
            return { labels: [], datasets: [], noDataMessage: 'No data available for this chart.' };
        }

        console.log('Incoming data for Victims by Corporation:', data); // Debugging log

        // Inspect each data item
        data.forEach((item, index) => {
            console.log(`Item ${index}:`, item);
        });

        // Update the field names based on the actual data structure
        const labels = data.map(item => item.name || 'Unknown');
        const victims = data.map(item => item.kill_count || 0);

        // Check for 'Unknown' labels
        const allUnknown = labels.every(label => label === 'Unknown');
        if (allUnknown) {
            console.warn(`All labels for ${chartName} are 'Unknown'. Check data source.`);
        }

        return {
            labels: labels,
            datasets: [{
                label: 'Number of Victims',
                data: victims,
                backgroundColor: 'rgba(255, 99, 132, 0.5)',
                borderColor: 'rgba(255, 99, 132, 1)',
                borderWidth: 1,
            }]
        };
    },
};

export default victimsByCorporationChartConfig;
