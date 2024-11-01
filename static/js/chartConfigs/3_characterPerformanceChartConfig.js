// static/js/chartConfigs/characterPerformanceChartConfig.js
import { truncateLabel, getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Character Performance Chart
 */
const characterPerformanceChartConfig = {
    id: 'characterPerformanceChart',
    instance: {}, // Object to hold instances per timeframe
    dataKeys: {
        mtd: { dataVar: 'mtdCharacterPerformanceData', canvasId: 'characterPerformanceChart_mtd' },
        ytd: { dataVar: 'ytdCharacterPerformanceData', canvasId: 'characterPerformanceChart_ytd' },
        lastMonth: { dataVar: 'lastMCharacterPerformanceData', canvasId: 'characterPerformanceChart_lastM' },
    },
    type: 'bar', // Base type
    options: getCommonOptions('Character Performance', {
        // No need to redefine plugins here as they're handled in getCommonOptions
        scales: {
            // y1 is already defined in getCommonOptions for Points
        },
        datasets: {
            // Define additional dataset options if needed
        },
    }),
    processData: function (data) {
        const chartName = 'Character Performance Chart';
        if (!validateChartDataArray(data, chartName)) {
            // Trigger noData plugin
            return { labels: [], datasets: [] };
        }

        // Sort data by KillCount descending
        const sortedData = [...data].sort((a, b) => (b.KillCount || 0) - (a.KillCount || 0));

        // Limit to top 10 characters
        const topN = 10;
        const limitedData = sortedData.slice(0, topN);

        const labels = limitedData.map(item => item.Name || 'Unknown'); // Ensure correct field name
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
                yAxisID: 'y', // Assign to primary y-axis
            },
            {
                label: 'Solo Kills',
                type: 'bar',
                data: soloKills,
                backgroundColor: 'rgba(153, 102, 255, 0.7)',
                borderColor: 'rgba(153, 102, 255, 1)',
                borderWidth: 1,
                yAxisID: 'y', // Assign to primary y-axis
            },
            {
                label: 'Points',
                type: 'line',
                data: points,
                backgroundColor: 'rgba(255, 159, 64, 0.7)',
                borderColor: 'rgba(255, 159, 64, 1)',
                borderWidth: 2,
                fill: false,
                yAxisID: 'y1', // Assign to secondary y-axis
                tension: 0.1,
                pointRadius: 4,
            },
        ];

        // Debugging Logs
        console.log('Processed Labels:', truncatedLabels);
        console.log('Kills Data:', kills);
        console.log('Solo Kills Data:', soloKills);
        console.log('Points Data:', points);
        console.log('Datasets:', datasets);

        return { labels: truncatedLabels, datasets, fullLabels: labels };
    },
};

export default characterPerformanceChartConfig;
