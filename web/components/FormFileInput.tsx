import React, {
  ChangeEvent,
  HTMLProps,
  ReactNode,
  useEffect,
  useRef,
} from "react";
import { useFormContext } from "react-hook-form";

export const FormFileInput = ({
  name,
  className,
  children,
  onChange,
  ...inputProps
}: {
  name: string;
  className?: string;
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
  children: (value: FileList) => ReactNode;
} & Omit<
  HTMLProps<HTMLInputElement>,
  "name" | "onChange" | "type" | "className" | "ref" | "children"
>) => {
  const { register, unregister, setValue, watch } = useFormContext();
  const inputRef = useRef<HTMLInputElement>(null);
  const value: FileList = watch(name);

  useEffect(() => {
    register(name);

    return () => unregister(name);
  }, [register, unregister, name]);

  return (
    <button
      className={className}
      onClick={() => inputRef.current?.click()}
      type="button"
    >
      <input
        ref={inputRef}
        className="hidden"
        type="file"
        onChange={(e) => {
          setValue(name, e.target.files);
          onChange?.(e);
        }}
        {...inputProps}
      />
      {children(value)}
    </button>
  );
};
