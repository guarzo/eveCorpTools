// static/js/chartConfigs/characterPerformanceChartConfig.js

import { truncateLabel, getColor, getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Character Performance Chart
 */
const characterPerformanceChartConfig = {
    type: 'bar', // Base type
    options: getCommonOptions('Character Performance', {
        scales: {
            // Additional scale options can be added here if needed
        },
        datasets: {
            // Additional dataset options can be added here if needed
        },
    }),
    processData: function (data) {
        const chartName = 'Character Performance Chart';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        console.log('Incoming data for Character Performance:', data); // Debugging log

        // Sort data by KillCount descending
        const sortedData = [...data].sort((a, b) => (b.KillCount || 0) - (a.KillCount || 0));

        // Limit to top 10 characters
        const topN = 10;
        const limitedData = sortedData.slice(0, topN);

        const labels = limitedData.map(item => item.CharacterName || item.Name || 'Unknown');
        const truncatedLabels = labels.map(label => truncateLabel(label, 15));

        const kills = limitedData.map(item => item.KillCount || 0);
        const soloKills = limitedData.map(item => item.SoloKills || 0);
        const points = limitedData.map(item => item.Points || 0);

        const datasets = [
            {
                label: 'Kills',
                type: 'bar',
                data: kills,
                backgroundColor: 'rgba(75, 192, 192, 0.7)',
                borderColor: 'rgba(75, 192, 192, 1)',
                borderWidth: 1,
                yAxisID: 'y',
            },
            {
                label: 'Solo Kills',
                type: 'bar',
                data: soloKills,
                backgroundColor: 'rgba(153, 102, 255, 0.7)',
                borderColor: 'rgba(153, 102, 255, 1)',
                borderWidth: 1,
                yAxisID: 'y',
            },
            {
                label: 'Points',
                type: 'line',
                data: points,
                backgroundColor: 'rgba(255, 159, 64, 0.7)',
                borderColor: 'rgba(255, 159, 64, 1)',
                borderWidth: 2,
                fill: false,
                yAxisID: 'y1',
                tension: 0.1,
                pointRadius: 4,
            },
        ];

        console.log('Processed Labels:', truncatedLabels);
        console.log('Kills Data:', kills);
        console.log('Solo Kills Data:', soloKills);
        console.log('Points Data:', points);
        console.log('Datasets:', datasets);

        return { labels: truncatedLabels, datasets, fullLabels: labels };
    },
};

export default characterPerformanceChartConfig;
