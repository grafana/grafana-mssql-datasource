// We set this specifically for 2 reasons.
// 1. It makes sense for both CI tests and local tests to behave the same so issues are found earlier
// 2. Any wrong timezone handling could be hidden if we use UTC/GMT local time (which would happen in CI).
process.env.TZ = 'Pacific/Easter'; // UTC-06:00 or UTC-05:00 depending on daylight savings

const { grafanaESModules, nodeModulesToTransform } = require('./.config/jest/utils');

module.exports = {
  ...require('./.config/jest.config'),
  // ESM-only packages pulled in by @grafana/* 13.3
  transformIgnorePatterns: [
    nodeModulesToTransform([...grafanaESModules, '@react-hookz/web', '@ver0/deep-equal', '@marcbachmann/cel-js']),
  ],
};
