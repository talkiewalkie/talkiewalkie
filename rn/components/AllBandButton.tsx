import { Text, TouchableOpacity, TouchableOpacityProps } from "react-native";
import React from "react";

const AllBandButton = ({
  label,
  style,
  ...props
}: { label: string } & TouchableOpacityProps) => (
  <TouchableOpacity
    style={{
      padding: 16,
      backgroundColor: "#002fa7",
      ...(style as object),
    }}
    {...props}
  >
    <Text style={{ color: "#fff", textAlign: "center" }}>{label}</Text>
  </TouchableOpacity>
);

export default AllBandButton;
