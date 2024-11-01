// static/js/chartConfigs/characterPerformanceChartConfig.js
import { truncateLabel, getColor, getCommonOptions, validateChartData } from '../utils.js';

const characterPerformanceChartConfig = {
    id: 'characterPerformanceChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdCharacterPerformanceData',
        ytd: 'ytdCharacterPerformanceData',
        lastMonth: 'lastMCharacterPerformanceData',
    },
    type: 'bar', // or any other suitable type
    options: getCommonOptions('Character Performance', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const dataPoint = context.raw; // Access the data point directly
                        const kills = dataPoint.Kills || 0;
                        const soloKills = dataPoint.SoloKills || 0;
                        const points = dataPoint.Points || 0;
                        return `Kills: ${kills}, Solo Kills: ${soloKills}, Points: ${points}`;
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'category',
                ticks: {
                    color: '#ffffff',
                    maxRotation: 45,
                    minRotation: 45,
                    autoSkip: false,
                },
                grid: { display: false },
                title: {
                    display: true,
                    text: 'Characters',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
            y: {
                beginAtZero: true,
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                title: {
                    display: true,
                    text: 'Performance Metrics',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
        },
    }),
    processData: function (data) {
        const chartName = 'Character Performance Chart';
        if (!validateChartData(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Sort data by Kills descending
        const sortedData = [...data].sort((a, b) => (b.Kills || 0) - (a.Kills || 0));

        // Limit to top 10 characters
        const topN = 10;
        const limitedData = sortedData.slice(0, topN);

        const labels = limitedData.map(item => item.CharacterName || 'Unknown');
        const truncatedLabels = labels.map(label => truncateLabel(label, 15));

        const kills = limitedData.map(item => item.Kills || 0);
        const soloKills = limitedData.map(item => item.SoloKills || 0);
        const points = limitedData.map(item => item.Points || 0);

        const datasets = [
            {
                label: 'Kills',
                data: kills,
                backgroundColor: 'rgba(75, 192, 192, 0.7)',
                borderColor: 'rgba(75, 192, 192, 1)',
                borderWidth: 1,
            },
            {
                label: 'Solo Kills',
                data: soloKills,
                backgroundColor: 'rgba(153, 102, 255, 0.7)',
                borderColor: 'rgba(153, 102, 255, 1)',
                borderWidth: 1,
            },
            {
                label: 'Points',
                data: points,
                backgroundColor: 'rgba(255, 159, 64, 0.7)',
                borderColor: 'rgba(255, 159, 64, 1)',
                borderWidth: 1,
            },
        ];

        return { labels: truncatedLabels, datasets, fullLabels: labels };
    },
};

export default characterPerformanceChartConfig;
