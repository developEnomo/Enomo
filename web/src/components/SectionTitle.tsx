import React from 'react';

type SectionTitleProps = {
  children: React.ReactNode;
};

const SectionTitle = ({ children }: SectionTitleProps) => {
  return <h2 className="text-lg font-bold text-[#861F6D]">{children}</h2>;
};

export default SectionTitle;