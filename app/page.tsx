import { BackGround } from "@/components/ui/squares-background";
import ShinyText from "@/components/ui/texts/shiny-text";
// import GradientText from "@/components/ui/texts/gradient-text";
import Navbar from "@/components/ui/navbar";

export default function SquaresDemo() {
  return (
    <div className="relative h-screen w-full">
      <BackGround />
      <Navbar />
      {
        /* Main content area */
        <div
          className="absolute top-20 left-0 w-full h-full flex items-start justify-center pt-32 pointer-events-none select-none"
          style={{ fontFamily: "SF-Pro-Display" }}
        >
          <div className="flex flex-col items-center gap-4">
            <div
              className="text-6xl font-bold text-white"
            >
              From seconds to sensations.
            </div>
            <ShinyText
              text="Short - Sharp - Shareable"
              className="font-thin"
              disabled={false}
              speed={3}
            />
          </div>
        </div>
      }
    </div>
  );
}
