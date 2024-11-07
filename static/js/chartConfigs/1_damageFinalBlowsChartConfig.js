// static/js/chartConfigs/1_damageFinalBlowsChartConfig.js

import { truncateLabel, getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Damage Done and Final Blows Chart
 */
const damageFinalBlowsChartConfig = {
    type: 'bar', // Base type for mixed charts
    options: getCommonOptions('Top Damage Done and Final Blows', {
        scales: {
            // Additional scale options can be added here if needed
        },
        datasets: {
            // Additional dataset options can be added here if needed
        },
        plugins: {
            tooltip: {
                mode: 'nearest', // Focus on the hovered bar segment
                intersect: true, // Show tooltip only when directly hovering over a segment
                callbacks: {
                    label: function (context) {
                        let labels = []; // Array to hold multiple dataset values
                        context.chart.data.datasets.forEach((dataset, index) => {
                            const label = dataset.label || '';
                            const value = dataset.data[context.dataIndex];
                            labels.push(`${label}: ${value.toLocaleString()}`);
                        });
                        return labels; //
                    },
                },
            },
        },
    }),
    processData: function (data) {
        const chartName = 'Damage Done and Final Blows';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Sort data by DamageDone descending
        const sortedData = [...data].sort((a, b) => (b.DamageDone || 0) - (a.DamageDone || 0));

        // Limit to top 20 characters
        const topN = 20;
        const limitedData = sortedData.slice(0, topN);

        const labels = limitedData.map(item => item.Name || 'Unknown');
        const truncatedLabels = labels.map(label => truncateLabel(label, 15)); // Truncate labels to 15 characters

        const damageData = limitedData.map(item => item.DamageDone || 0);
        const finalBlowsData = limitedData.map(item => item.FinalBlows || 0);

        const datasets = [
            {
                label: 'Damage Done',
                type: 'bar', // Explicitly set type as bar
                data: damageData,
                backgroundColor: 'rgba(255, 77, 77, 0.7)',
                borderColor: 'rgba(255, 77, 77, 1)',
                borderWidth: 1,
                yAxisID: 'y', // Assign to primary y-axis
            },
            {
                label: 'Final Blows',
                type: 'line', // Set type as line
                data: finalBlowsData,
                backgroundColor: 'rgba(54, 162, 235, 0.7)',
                borderColor: 'rgba(54, 162, 235, 1)',
                borderWidth: 2,
                fill: false,
                yAxisID: 'y1', // Assign to secondary y-axis
                tension: 0.1, // Smoothness of the line
                pointRadius: 4, // Size of points on the line
            },
        ];

        return { labels: truncatedLabels, datasets, fullLabels: labels };
    },
};

export default damageFinalBlowsChartConfig;
