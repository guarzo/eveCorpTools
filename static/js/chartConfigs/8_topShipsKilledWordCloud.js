// static/js/chartConfigs/8_topShipsKilledWordCloudConfig.js

import { validateChartDataArray } from '../utils.js';
import { wordCloudOptions } from '../wordCloudOptions.js'; // Import the separate options

/**
 * Configuration for the Top Ships Killed Word Cloud
 */
const topShipsKilledWordCloudConfig = {
    id: 'topShipsKilledWordCloud',
    instance: {}, // Object to track instances by timeframe
    dataKeys: {
        mtd: { dataVar: 'mtdTopShipsKilledData', canvasId: 'topShipsKilledWordCloud_mtd' },
        ytd: { dataVar: 'ytdTopShipsKilledData', canvasId: 'topShipsKilledWordCloud_ytd' },
        lastMonth: { dataVar: 'lastMTopShipsKilledData', canvasId: 'topShipsKilledWordCloud_lastM' },
    },
    type: 'wordCloud',
    dataType: 'array',
    options: wordCloudOptions, // Use the separate Word Cloud options

    processData: function (data) {
        const chartName = 'Top Ships Killed Word Cloud';

        // Validate and log data
        const isValidData = validateChartDataArray(data, chartName);
        console.log(`Data validation for ${chartName}:`, isValidData);

        if (!isValidData) {
            return { labels: [], datasets: [] };
        }

        // Filter and transform data
        const filteredData = data.filter(item => item.KillCount > 0);
        const wordCloudData = filteredData.map(item => ({
            text: item.Name || 'Unknown',
            weight: item.KillCount || 1,
        }));

        // Debugging: Log the transformed data
        console.log(`Processed data for ${chartName}:`, wordCloudData);

        // Return structured data for Word Cloud chart
        return {
            labels: [],  // Word Cloud doesn't use labels
            datasets: [{
                label: 'Killmails',
                data: wordCloudData,
                // Remove rotation and color for testing
            }]
        };
    },
};

export default topShipsKilledWordCloudConfig;
