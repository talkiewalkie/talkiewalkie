import { useRef } from "react";
import { Animated, PanResponder, View, Text, StyleSheet } from "react-native";

// unused component, drawn from: https://reactnative.dev/docs/panresponder
// other interesting code snippets for the usecase:
// - https://gist.github.com/himanshu717171/b70794816b36b1de41d8dd679e17aba5
// - https://reactnative.dev/docs/panresponder

export const Draggable = () => {
  const pan = useRef(new Animated.ValueXY()).current;

  const panResponder = useRef(
    PanResponder.create({
      onMoveShouldSetPanResponder: () => true,
      onPanResponderGrant: () => {
        pan.setOffset({
          x: (pan.x as any)._value,
          y: (pan.y as any)._value,
        });
      },
      onPanResponderMove: Animated.event([null, { dx: pan.x, dy: pan.y }]),
      onPanResponderRelease: () => {
        pan.flattenOffset();
      },
    })
  ).current;

  return (
    <View
      style={{
        flex: 1,
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Text
        style={{
          fontSize: 14,
          lineHeight: 24,
          fontWeight: "bold",
        }}
      >
        Drag this box!
      </Text>
      <Animated.View
        style={{
          transform: [{ translateX: pan.x }, { translateY: pan.y }],
        }}
        {...panResponder.panHandlers}
      >
        <View
          style={{
            height: 150,
            width: 150,
            backgroundColor: "blue",
            borderRadius: 5,
          }}
        />
      </Animated.View>
    </View>
  );
};
