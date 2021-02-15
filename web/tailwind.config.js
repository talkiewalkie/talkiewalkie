const config = require("./src/styles/config.json");

// List of all available plugins: https://tailwindcss.com/docs/configuration#core-plugins
// For reference, the default configuration: https://github.com/tailwindcss/tailwindcss/blob/master/stubs/defaultConfig.stub.js
module.exports = {
  purge: ["./src/**/*.ts", "./src/**/*.tsx"],
  future: { removeDeprecatedGapUtilities: true },
  theme: {
    screens: {
      lg: "1024px",
    },
    colors: config.colors,
    textColor: (theme) => ({ ...theme("colors"), placeholder: "#c9d0e1" }),
    backgroundColor: (theme) => ({
      ...theme("colors"),
      "radical-red": "#ff266E",
      "primary-100": "rgba(0, 19, 255, 0.1)",
      "primary-50": "rgba(0, 19, 255, 0.05)",
      overlay: "rgba(31, 51, 86, 0.8)",
      "white-overlay": "rgba(255, 255, 255, 0.95)",
      "white-overlay-800": "rgba(255, 255, 255, 0.8)",
      "white-overlay-300": "rgba(255, 255, 255, 0.3)",
      "blue-overlay": "rgba(245, 251, 255, 0.95)",
      "light-green-100": "rgba(174, 207, 214, 0.1)",
      "coral-100": "rgba(239, 187, 172, 0.1)",
      "danger-100": "rgba(250, 74, 72, 0.1)",
    }),

    borderColor: (theme) => ({
      transparent: theme("colors.transparent"),
      white: theme("colors.white"),
      primary: theme("colors.primary"),
      danger: theme("colors.danger"),
      stroke: theme("colors.stroke"),
      default: theme("colors.stroke"),
      orange: theme("colors.orange"),
      "primary-dark": theme("colors.primary-dark"),
      coral: theme("colors.coral"),
      "light-green": theme("colors.light-green"),
      "primary-100": "rgba(0, 19, 255, 0.1)",
    }),
    borderRadius: {
      0: "0",
      sm: "5px",
      default: "10px",
      full: "9999px",
    },
    borderWidth: {
      0: "0",
      default: "1px",
      4: "4px",
    },
    boxShadow: {
      none: "none",
      xs: "0 1px 2px 0 rgba(38, 40, 45, 0.05)",
      sm: "0 5px 10px 0 rgba(0, 82, 255, 0.05)",
      "sm-outlined":
        "0 4px 8px 0 rgba(0, 19, 255, 0.05), 0 1px 2px 0 rgba(0, 19, 255, 0.08)",
      lg: "0 10px 20px rgba(0,0,0,0.08)",
    },
    divideColor: (theme) => ({
      primary: theme("colors.primary"),
      default: theme("colors.stroke"),
    }),
    divideWidth: {
      default: "1px",
    },

    fontSize: {
      10: "10px",
      11: "11px",
      12: "12px",
      13: "13px",
      15: "15px",
      18: "18px",
      20: "20px",
      24: "24px",
      30: "30px",
      36: "36px",
    },
    fontWeight: {
      normal: "400",
      medium: "500",
      semibold: "600",
    },

    spacing: {
      0: "0",
      4: "4px",
      6: "6px",
      8: "8px",
      10: "10px",
      12: "12px",
      14: "14px",
      16: "16px",
      18: "18px",
      20: "20px",
      24: "24px",
      30: "30px",
      32: "32px",
      36: "36px",
      44: "44px",
      52: "52px",
      56: "56px",
    },
    margin: (theme, { negative }) => ({
      ...theme("spacing"),
      ...negative(theme("spacing")),
      auto: "auto",
    }),
    padding: (theme) => theme("spacing"),
    height: (theme) => ({
      ...theme("spacing"),
      1: "1px",
      2: "2px",
      auto: "auto",
      full: "100%",
    }),
    width: (theme) => ({
      ...theme("spacing"),
      auto: "auto",
      "1/2": "50%",
      "1/3": "33.33%",
      "1/4": "25%",
      full: "100%",
    }),
    maxHeight: {
      0: "0",
      full: "100%",
      screen: "100vh",
    },
    maxWidth: {
      "1/2": "50%",
      330: "330px",
      440: "440px",
      full: "100%",
      screen: "100vw",
    },

    flex: {
      1: "1",
      2: "2",
    },
    zIndex: config.zIndex,
    inset: (theme, { negative }) => ({
      ...theme("spacing"),
      ...negative(theme("spacing")),
      0: "0",
      full: "100%",
      initial: "initial",
    }),
    opacity: {
      0: "0", // Weird: quote are needed to not be set as px
      10: "0.1",
      50: "0.5",
      70: "0.7",
      80: "0.8",
      90: "0.9",
      100: "1",
    },

    transitionProperty: {
      default: "all",
      none: "none",
    },
    transitionDuration: {
      100: "100ms",
      200: "200ms",
      300: "300ms",
    },
    rotate: {
      "-90": "-90deg",
      "-45": "-45deg",
      45: "45deg",
      90: "90deg",
      135: "135deg",
      180: "180deg",
    },
    translate: {
      full: "100%",
      "-1/2": "-50%",
      "-full": "-100%",
    },
  },
  variants: {
    alignContent: [],
    alignItems: [],
    alignSelf: [],
    appearance: [],
    backgroundAttachment: [],
    backgroundColor: ["hover", "active", "disabled"],
    backgroundOpacity: [],
    borderCollapse: [],
    borderColor: ["hover", "active"],
    borderRadius: [],
    borderWidth: ["first", "last", "responsive"],
    boxShadow: ["hover", "focus", "responsive"],
    cursor: [],
    display: ["responsive", "hover", "group-hover"],
    divideColor: [],
    divideWidth: [],
    flex: [],
    flexGrow: [],
    flexShrink: [],
    flexWrap: [],
    fontSize: ["responsive"],
    fontStyle: [],
    fontWeight: ["hover"],
    height: ["responsive"],
    inset: ["responsive"],
    justifyContent: [],
    margin: ["responsive", "first", "last"],
    space: ["responsive"],
    maxHeight: ["responsive"],
    maxWidth: ["responsive"],
    minHeight: [],
    minWidth: [],
    objectFit: [],
    opacity: ["hover", "group-hover", "active", "disabled"],
    outline: ["focus"],
    padding: ["responsive", "first", "last"],
    pointerEvents: [],
    position: ["responsive"],
    resize: [],
    textAlign: ["responsive"],
    textColor: ["hover", "group-hover", "active"],
    textDecoration: ["hover", "group-hover", "active"],
    textTransform: [],
    verticalAlign: [],
    whitespace: [],
    width: ["responsive"],
    wordBreak: [],
    visibility: ["responsive", "hover", "group-hover"],
    zIndex: [],
    transform: [],
    rotate: [],
    transitionProperty: [],
    transitionTimingFunction: [],
    transitionDuration: [],
  },
  corePlugins: {
    accessibility: false,
    backgroundPosition: false,
    backgroundRepeat: false,
    backgroundSize: false,
    borderOpacity: false,
    borderStyle: true,
    boxSizing: false,
    container: false,
    divideOpacity: false,
    fill: false,
    float: false,
    clear: false,
    fontFamily: false,
    fontSmoothing: false,
    lineHeight: false,
    letterSpacing: false,
    objectPosition: false,
    order: false,
    placeholderColor: false,
    placeholderOpacity: false,
    stroke: false,
    strokeWidth: false,
    tableLayout: false,
    textOpacity: false,
    userSelect: false,
    flexDirection: false, // Custom version
    flexWrap: false, // Custom version
    overflow: false, // Custom version
    gap: true,
    gridTemplateColumns: true,
    gridColumn: true,
    gridColumnStart: true,
    gridColumnEnd: true,
    gridTemplateRows: true,
    gridRow: true,
    gridRowStart: true,
    gridRowEnd: true,
    gridAutoFlow: true,
    transformOrigin: false,
    scale: false,
    space: false, // Custom version
    translate: true,
    skew: false,
    transitionDelay: false,
  },
  plugins: [
    // The default space plugin cannot be override by a child margin
    // This is a custom version with less specification
    ({ addComponents, theme, e }) => {
      const space = Object.entries(theme("space"))
        .filter(([k]) => !k.startsWith("-")) // Don't generate margin version
        .reduce(
          (acc, [k, v]) => ({
            ...acc,
            [`.${e(`space-x-${k}`)} > * + *`]: { marginLeft: v },
            [`.${e(`space-y-${k}`)} > * + *`]: { marginTop: v },
          }),
          {}
        );
      // The problem is that custom utilities are injected at the end
      // Will be resolved with this feature request: https://github.com/tailwindcss/tailwindcss/issues/1069
      // For now, the hack is to use addComponents, with are injected before utilities
      // The problem is that variant version are not generated. Ok for now
      addComponents(space);
    },
    ({ addComponents }) => {
      addComponents({
        ".label": {
          letterSpacing: "0.2px",
          fontWeight: 600,
          textTransform: "uppercase",
        },
      });
    },
    // Custom version for flexDirection & flexWrap to add `display: flex` by default on `flex-col` & `flex-wrap`
    ({ addUtilities }) => {
      addUtilities(
        {
          ".flex-row": { flexDirection: "row" },
          ".flex-row-reverse": {
            display: "flex",
            flexDirection: "row-reverse",
          },
          ".flex-col": { display: "flex", flexDirection: "column" },
          ".flex-col-reverse": {
            display: "flex",
            flexDirection: "column-reverse",
          },
          ".flex-wrap": { display: "flex", flexWrap: "wrap" },
          ".flex-wrap-reverse": { display: "flex", flexWrap: "wrap-reverse" },
          ".flex-no-wrap": { flexWrap: "nowrap" },
          // Also add more complex but useful ones
          ".flex-center": {
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          },
          // This class should be used on the "main" element of a row or a column
          // The goal is to take the place available after some fixed element sizes (header, icon, time, ...)
          // If you want multiple main elements to share the available space, use flex-1 or flex-2 or flex-1/2
          // This SO answer explain very well the need for min-width: https://stackoverflow.com/a/36247448/6851875
          ".flex-fill": {
            flex: "1 1 auto",
            minWidth: 0,
            minHeight: 0,
          },
          // This is to force two elements to share available space equally
          ".flex-1\\/2": {
            flex: "1 1 50%",
          },
        },
        ["responsive"]
      );
    },
    // Avoid creating visible & scroll versions, that are 99% of the time not desired
    ({ addUtilities }) => {
      addUtilities(
        {
          ".overflow-auto": { overflow: "auto" },
          ".overflow-x-auto": { "overflow-x": "auto" },
          ".overflow-y-auto": { "overflow-y": "auto" },
          ".overflow-hidden": { overflow: "hidden" },
          ".overflow-x-hidden": { "overflow-x": "hidden" },
          ".overflow-y-hidden": { "overflow-y": "hidden" },
        },
        ["responsive"]
      );
    },
  ],
};
