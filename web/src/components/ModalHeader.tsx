"use client";

import { useRouter } from 'next/navigation';
import { FaTimes } from 'react-icons/fa';

const ModalHeader = () => {
  const router = useRouter();

  return (
    <header className="flex justify-end items-center p-4">
      <button
        onClick={() => router.back()}
        className="text-4xl text-[#C73BA4]"
        aria-label="Close"
      >
        <FaTimes />
      </button>
    </header>
  );
};

export default ModalHeader;