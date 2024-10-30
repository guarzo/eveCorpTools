// static/js/chartConfigs/killActivityChartConfig.js

import { getCommonOptions } from '../utils.js';

/**
 * Configuration for the Kill Activity Over Time Chart
 */
const killActivityChartConfig = {
    id: 'killActivityChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdKillActivityData',
        ytd: 'ytdKillActivityData',
        lastMonth: 'lastMKillActivityData',
    },
    type: 'line',
    options: getCommonOptions('Kill Activity Over Time', {
        plugins: {
            legend: { display: false },
        },
        scales: {
            x: {
                type: 'time',
                time: {
                    unit: 'day',
                },
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
            },
            y: {
                beginAtZero: true,
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
            },
        },
    }),
    processData: function (data) {
        const labels = data.map(item => new Date(item.Time));
        const kills = data.map(item => item.Value);

        const datasets = [{
            label: 'Kills Over Time',
            data: kills,
            borderColor: 'rgba(255, 77, 77, 1)',
            backgroundColor: 'rgba(255, 77, 77, 0.5)',
            fill: true,
        }];

        return { labels, datasets };
    },
};

export default killActivityChartConfig;
