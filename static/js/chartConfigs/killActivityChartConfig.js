// chartConfigs/killActivityChartConfig.js
// import { commonOptions } from '../utils.js';

import {commonOptions} from "../utils.js";

const killActivityChartConfig = {
    id: 'killActivityChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdKillActivityData',
        ytd: 'ytdKillActivityData',
        lastMonth: 'lastMKillActivityData',
    },
    type: 'line',
    options: {
        ...commonOptions,
        plugins: {
            ...commonOptions.plugins,
            legend: { display: false },
        },
        scales: {
            x: {
                ...commonOptions.scales.x,
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
    },
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

//
// // chartConfigs/killActivityChartConfig.js
// import { commonOptions } from '../utils.js';
//
// const killActivityChartConfig = {
//     id: 'killActivityChart',
//     instance: null,
//     dataKeys: {
//         mtd: 'mtdKillActivityData',
//         ytd: 'ytdKillActivityData',
//         lastMonth: 'lastMKillActivityData',
//     },
//     type: 'line',
//     options: {
//         ...commonOptions,
//         plugins: {
//             ...commonOptions.plugins,
//             legend: { display: false },
//         },
//         scales: {
//             x: {
//                 ...commonOptions.scales.x,
//                 type: 'time',
//                 time: {
//                     unit: 'day',
//                 },
//             },
//             y: {
//                 ...commonOptions.scales.y,
//             },
//         },
//     },
//     processData: function (data) {
//         // ... same as before ...
//     },
// };
//
// export default killActivityChartConfig;
