// static/js/chartConfigs/killHeatmapChartConfig.js

import {getColor, getCommonOptions, validateChartData} from '../utils.js';

/**
 * Configuration for the Kill Heatmap Chart
 */
const killHeatmapChartConfig = {
    id: 'killHeatmapChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdKillHeatmapData',
        ytd: 'ytdKillHeatmapData',
        lastMonth: 'lastMKillHeatmapData',
    },
    type: 'matrix',
    options: getCommonOptions('Kill Heatmap', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const x = context.raw.x;
                        const y = context.raw.y;
                        const value = context.raw.v;
                        const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
                        return `Day: ${days[y]}, Hour: ${x}, Kills: ${value}`;
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'category',
                labels: [...Array(24).keys()],
                ticks: { color: '#ffffff' },
                grid: { display: false },
            },
            y: {
                type: 'category',
                labels: ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'],
                ticks: { color: '#ffffff' },
                grid: { display: false },
            },
        },
    }),
    processData: function (data) {
        const chartName = 'Kill Heatmap';
        if (!validateChartData(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }
        const maxKills = Math.max(...data.flat());
        const heatmapData = [];

        for (let day = 0; day < 7; day++) {
            for (let hour = 0; hour < 24; hour++) {
                const kills = data[day][hour];
                heatmapData.push({
                    x: hour,
                    y: day,
                    v: kills,
                });
            }
        }

        const datasets = [{
            label: 'Kill Heatmap',
            data: heatmapData,
            backgroundColor: function (ctx) {
                const value = ctx.dataset.data[ctx.dataIndex].v;
                const alpha = maxKills ? value / maxKills : 0;
                return `rgba(255, 77, 77, ${alpha})`;
            },
            width: ({ chart }) => (chart.chartArea || {}).width / 24 - 1,
            height: ({ chart }) => (chart.chartArea || {}).height / 7 - 1,
        }];

        return { datasets };
    },
};

export default killHeatmapChartConfig;
