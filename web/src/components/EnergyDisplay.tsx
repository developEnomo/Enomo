"use client"
import Image from "next/image"
import { useEffect, useRef } from "react"

const energyImageMap: { [key: number]: string } = {
  1: "/Images/EnergyIcons/Energy_full_1.png",
  2: "/Images/EnergyIcons/Energy_full_2.png",
  3: "/Images/EnergyIcons/Energy_full_3.png",
  4: "/Images/EnergyIcons/Energy_full_4.png",
}

type EnergyDisplayProps = {
  energyLevels: number[]
}

export default function EnergyDisplay({ energyLevels }: EnergyDisplayProps) {
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const scrollContainer = scrollRef.current
    if (!scrollContainer) return

    const scrollSpeed = 0.5
    let animationFrameId: number

    const animateScroll = () => {
      scrollContainer.scrollLeft += scrollSpeed

      const maxScrollLeft = scrollContainer.scrollWidth - scrollContainer.clientWidth
      if (scrollContainer.scrollLeft >= maxScrollLeft) {
        scrollContainer.scrollLeft = 0
      }

      animationFrameId = requestAnimationFrame(animateScroll)
    }

    animationFrameId = requestAnimationFrame(animateScroll)

    return () => cancelAnimationFrame(animationFrameId)
  }, [])

  // 表示する配列：元データ ×2 ＋ 空白 ＋ 先頭4個
  const duplicatedLevels = [
       "gap",
    ...energyLevels,
    "gap", // 空白の印
  ]

  return (
    <div className="w-full">
      <h2 className="text-lg font-bold text-center mb-4 text-gray-800">
        みんなのエナジー
      </h2>

      <div
        ref={scrollRef}
        className="flex items-center space-x-4 p-2 border-2 border-gray-200 rounded-lg overflow-x-auto flex-nowrap scrollbar-hide"
      >
        {duplicatedLevels.map((item, index) => {
          if (item === "gap") {
            return <div key={index} className="flex-shrink-0 w-full" />
          }

          const imageUrl =
            energyImageMap[item as number] || "/Images/EnergyIcons/Energy_empty.svg"
          return (
            <div key={index} className="flex-shrink-0">
              <Image
                src={imageUrl}
                alt={`エナジーレベル ${item}`}
                width={60}
                height={60}
                priority={index < 5}
              />
            </div>
          )
        })}
      </div>
    </div>
  )
}