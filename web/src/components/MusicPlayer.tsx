type MusicPlayerProps = {
  trackId: string;
};

export default function MusicPlayer({ trackId }: MusicPlayerProps) {
  const spotifyEmbedUrl = `https://open.spotify.com/embed/track/${trackId}`;

  return (
    <div className="w-full">
      <iframe
        style={{ borderRadius: '12px' }}
        src={spotifyEmbedUrl}
        width="100%"
        height="152"
        frameBorder="0"
        allowFullScreen={false}
        allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture"
        loading="lazy"
      ></iframe>
    </div>
  );
}