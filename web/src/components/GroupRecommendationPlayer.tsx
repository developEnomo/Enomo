"use client";

import { useEffect, useState } from "react";
import MusicPlayer from "./MusicPlayer";

type Track = { id: string };
type RecommendationsResponse = { tracks: Track[] };

export default function GroupRecommendationPlayer({ groupId }: { groupId: string }) {
  const [trackId, setTrackId] = useState<string | null>(null);
  const [status, setStatus] = useState<"loading" | "ok" | "error">("loading");

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const res = await fetch(`/api/v1/groups/${encodeURIComponent(groupId)}/recommendations?limit=1`, {
          credentials: "include", // Cookie セッション前提
          cache: "no-store",
        });
        if (!res.ok) throw new Error(`status ${res.status}`);
        const data: RecommendationsResponse = await res.json();
        if (cancelled) return;

        const id = data?.tracks?.[0]?.id ?? null;
        setTrackId(id);
        setStatus("ok");
      } catch {
        if (!cancelled) setStatus("error");
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [groupId]);

  if (status === "loading") return <div>Loading...</div>;
  if (status === "error" || !trackId) return <div>曲が見つかりませんでした</div>;

  return <MusicPlayer trackId={trackId} />;
}
