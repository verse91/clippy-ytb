import { BackGround } from "@/components/ui/squares-background";
import ShinyText from "@/components/ui/texts/shiny-text";
import { TextShimmerWave } from "@/components/ui/texts/text-wave";
import Navbar from "@/components/ui/navbar";
import { BoxChat } from "@/components/ui/box-chat";

export default function Home() {
  return (
    <div className="relative min-h-screen w-full">
      <div className="fixed inset-0 z-0">
        <BackGround />
      </div>
      <div className="relative z-10">
        <Navbar />
        <div
          className="w-full h-full flex items-start justify-center pt-48"
          style={{ fontFamily: "SF-Pro-Display" }}
        >
          <div className="flex flex-col items-center gap-4">
            <TextShimmerWave
              className="text-6xl font-bold pointer-events-none select-none"
              duration={3}
            >
              From seconds to sensations.
            </TextShimmerWave>
            <ShinyText
              text="Short - Sharp - Shareable"
              className="font-thin pointer-events-none select-none"
              disabled={false}
              speed={3}
            />
            <div className="flex w-screen overflow-x-hidden -mt-48">
              <BoxChat />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
