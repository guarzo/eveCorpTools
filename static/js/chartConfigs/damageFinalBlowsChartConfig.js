// chartConfigs/damageFinalBlowsChartConfig.js
import { truncateLabel } from '../utils.js';

const damageFinalBlowsChartConfig = {
    id: 'damageFinalBlowsChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdCharacterDamageData',
        ytd: 'ytdCharacterDamageData',
        lastMonth: 'lastMCharacterDamageData',
    },
    type: 'bar',
    options: {
        indexAxis: 'y',
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: { display: true },
            tooltip: {
                callbacks: {
                    title: function (context) {
                        const index = context[0].dataIndex;
                        return context[0].chart.data.fullLabels
                            ? context[0].chart.data.fullLabels[index]
                            : context[0].chart.data.labels[index];
                    },
                },
            },
        },
        scales: {
            x: {
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                beginAtZero: true,
            },
            y: {
                ticks: {
                    color: '#ffffff',
                    autoSkip: false,
                },
                grid: { display: false },
            },
        },
    },
    processData: function (data) {
        const labels = data.map(item => item.Name || 'Unknown');
        const damageData = data.map(item => item.DamageDone || 0);
        const finalBlowsData = data.map(item => item.FinalBlows || 0);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const datasets = [
            {
                label: 'Damage Done',
                data: damageData,
                backgroundColor: 'rgba(255, 77, 77, 0.7)',
            },
            {
                label: 'Final Blows',
                data: finalBlowsData,
                backgroundColor: 'rgba(54, 162, 235, 0.7)',
            },
        ];
        return { labels: truncatedLabels, datasets, fullLabels };
    },
};

export default damageFinalBlowsChartConfig;



// // chartConfigs/damageFinalBlowsChartConfig.js
// import { truncateLabel, commonOptions } from '../utils.js';
//
// const damageFinalBlowsChartConfig = {
//     id: 'damageFinalBlowsChart',
//     instance: null,
//     dataKeys: {
//         mtd: 'mtdCharacterDamageData',
//         ytd: 'ytdCharacterDamageData',
//         lastMonth: 'lastMCharacterDamageData',
//     },
//     type: 'bar',
//     options: {
//         ...commonOptions,
//         indexAxis: 'y',
//         scales: {
//             ...commonOptions.scales,
//             x: {
//                 ...commonOptions.scales.x,
//                 beginAtZero: true,
//             },
//             y: {
//                 ...commonOptions.scales.y,
//                 ticks: {
//                     ...commonOptions.scales.y.ticks,
//                     autoSkip: false,
//                 },
//                 grid: { display: false },
//             },
//         },
//     },
//     processData: function (data) {
//         // ... same as before ...
//     },
// };
//
// export default damageFinalBlowsChartConfig;
