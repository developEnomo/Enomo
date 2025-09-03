import Image from 'next/image';

const AppLogo = () => {
    return (
        <div className="flex justify-center py-8">
            <Image
                src="/Images/Iogos/Logo_Enomo_Icon.svg"
                alt="Enomo App Icon"
                width={80}
                height={80}
            />
        </div>
    );
};

export default AppLogo;