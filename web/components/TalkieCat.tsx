import cat from "../public/gifs/cat.gif";
import Image from "next/image";

export default function TalkieCat({
  fixed = false,
  className,
}: {
  fixed?: boolean;
  className?: string;
}) {
  return fixed ? (
    <div className={className}>
      <Image src={cat} alt="TalkieWalkie cat" />
    </div>
  ) : (
    <div
      className={className}
      style={{
        position: "fixed",
        bottom: 120,
        right: -150,
        animation: "8s linear infinite",
        animationDelay: "3s",
        animationName: "slidein",
      }}
    >
      <Image src={cat} alt="TalkieWalkie cat" />
    </div>
  );
}
