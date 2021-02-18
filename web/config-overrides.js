const { override, addPostcssPlugins } = require("customize-cra");

module.exports = override(
  addPostcssPlugins([
    require("postcss-simple-vars"),
    require("tailwindcss"),
    require("postcss-hexrgba"),
    require("postcss-nested"),
    require("autoprefixer"),
  ]),
);
