// static/js/chartConfigs/characterPerformanceChartConfig.js

import { truncateLabel, getColor, getCommonOptions, validateChartData } from '../utils.js';

/**
 * Configuration for the Character Performance Chart
 */
const characterPerformanceChartConfig = {
    id: 'characterPerformanceChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdCharacterPerformanceData',
        ytd: 'ytdCharacterPerformanceData',
        lastMonth: 'lastMCharacterPerformanceData',
    },
    type: 'bar',
    options: getCommonOptions('Character Performance', {
        indexAxis: 'y', // Switch to horizontal bar chart
        plugins: {
            legend: { display: true, position: 'top', labels: { color: '#ffffff' } },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const label = context.dataset.label || '';
                        const value = context.parsed.y !== undefined ? context.parsed.y : context.parsed.x;
                        return `${label}: ${value}`;
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'linear',
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                beginAtZero: true,
                title: {
                    display: true,
                    text: 'Count',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
            y: {
                type: 'category',
                labels: [], // Labels are set dynamically in processData
                ticks: {
                    color: '#ffffff',
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
        },
    }),
    processData: function (data) {
        const chartName = 'Character Performance';
        if (!validateChartData(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Extract labels and data
        const labels = data.map(item => item.CharacterName || 'Unknown');
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const killCountData = data.map(item => item.KillCount || 0);
        const soloKillsData = data.map(item => item.SoloKills || 0);
        const pointsData = data.map(item => item.Points || 0);

        // Define datasets
        const datasets = [
            {
                label: 'Kill Count',
                data: killCountData,
                backgroundColor: 'rgba(255, 77, 77, 0.7)',
                borderColor: 'rgba(255, 77, 77, 1)',
                borderWidth: 1,
            },
            {
                label: 'Solo Kills',
                data: soloKillsData,
                backgroundColor: 'rgba(54, 162, 235, 0.7)',
                borderColor: 'rgba(54, 162, 235, 1)',
                borderWidth: 1,
            },
            {
                label: 'Points',
                data: pointsData,
                backgroundColor: 'rgba(255, 206, 86, 0.7)',
                borderColor: 'rgba(255, 206, 86, 1)',
                borderWidth: 1,
            },
        ];

        return { labels: truncatedLabels, datasets, fullLabels: labels };
    },
};

export default characterPerformanceChartConfig;
