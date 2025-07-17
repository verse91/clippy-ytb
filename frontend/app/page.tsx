import { BackGround } from "@/components/squares-background";
import ShinyText from "@/components/ui/texts/shiny-text";
import { TextShimmerWave } from "@/components/ui/texts/text-wave";
import Navbar from "@/components/navbar";
import { BoxChat } from "@/components/box-chat";

export default function Home() {
  return (
    <div className="relative min-h-screen w-full">
      <div className="fixed inset-0 z-0">
        <BackGround />
      </div>

      <div className="relative z-10">
        <Navbar />
        <div
          className="w-full h-full flex items-start justify-center pt-16 sm:pt-24 md:pt-32 lg:pt-48 px-4 sm:px-6 lg:px-8"
          style={{ fontFamily: "SF-Pro-Display" }}
        >
          <div className="flex flex-col items-center gap-2 sm:gap-4 w-full max-w-7xl">
            <TextShimmerWave
              className="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-bold pointer-events-none select-none text-center px-4"
              duration={3}
            >
              From seconds to sensations.
            </TextShimmerWave>
            <div className="flex flex-col items-center w-full gap-4 sm:gap-6 lg:gap-10">
              <ShinyText
                text="Turn raw clips into viral hits with Clippy"
                className="font-thin pointer-events-none select-none text-center px-4 text-sm sm:text-base md:text-lg"
                disabled={false}
                speed={3}
              />
              <div className="flex w-full max-w-4xl overflow-x-hidden px-4 sm:px-6 lg:px-8">
                <BoxChat />
                <div className="min-h-screen absolute inset-0 w-full h-full overflow-hidden pointer-events-none z-0">
                  <div className="absolute top-0 left-1/4 w-96 h-96 bg-violet-500/10 rounded-full mix-blend-normal filter blur-[128px] animate-pulse" />
                  <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-indigo-500/10 rounded-full mix-blend-normal filter blur-[128px] animate-pulse delay-700" />
                  <div className="absolute top-1/4 right-1/3 w-64 h-64 bg-fuchsia-500/10 rounded-full mix-blend-normal filter blur-[96px] animate-pulse delay-1000" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
