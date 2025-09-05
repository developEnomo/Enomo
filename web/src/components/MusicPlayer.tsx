"use client";

type MusicPlayerProps = {
  trackId: string;
};

export default function MusicPlayer({ trackId }: MusicPlayerProps) {
  const src = `https://open.spotify.com/embed/track/${encodeURIComponent(trackId)}`;

  return (
    <div className="w-full">
      <iframe
        key={src} // trackId が変わったら確実にリロード
        title="Spotify Music Player"
        style={{ borderRadius: 12 }}
        src={src}
        width="100%"
        height={152}
        frameBorder="0"
        allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture"
        allowFullScreen={false}
        loading="lazy"
      />
    </div>
  );
}
