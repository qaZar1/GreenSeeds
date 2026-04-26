import React, { useEffect, useState } from "react";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import { StatCard } from "../../utils/Card";
import toast from "react-hot-toast";

const DashboardPage: React.FC = () => {

  usePageHeader(
    "Дашборд",
    "Общее состояние системы"
  );

  const [seeds, setSeeds] = useState<any[]>([]);
  const [bunkers, setBunkers] = useState<any[]>([]);
  const [placements, setPlacements] = useState<any[]>([]);
  const [devices, setDevices] = useState<any[]>([]);

  useEffect(() => {

    const load = async () => {

      try {

        const [seedsData, bunkersData, placementsData, devicesData] = await Promise.all([
          api.getList("seeds"),
          api.getList("bunkers"),
          api.getList("placements"),
          api.getList("devices")
        ]);

        setSeeds(seedsData || []);
        setBunkers(bunkersData || []);
        setPlacements(placementsData || []);
        setDevices(devicesData || []);

      } catch {

        toast.error("Ошибка загрузки данных");

      }

    };

    load();

  }, []);

  const activeSeeds = seeds.filter(s => !s.deleted_at);

  const emptyBunkers =
    bunkers.length - placements.length;

  const activePlacements =
    placements.length;

  const onlineDevices =
    devices.filter(d => d.online).length;

  const offlineDevices =
    devices.length - onlineDevices;

  return (

    <div className="grid grid-cols-1 lg:grid-cols-3 gap-[16px]">

  <StatCard
    title="Всего семян"
    value={seeds.length}
  />

  <StatCard
    title="Активные семена"
    value={activeSeeds.length}
    color="var(--status-success-text)"
  />

  <StatCard
    title="Пустые бункеры"
    value={emptyBunkers}
    color="var(--status-warning-text)"
  />

  {/* большой блок */}

  <div className="lg:col-span-2 bg-[var(--bg-secondary)] rounded-[12px] p-[20px]">

    <div className="text-[14px] text-[var(--text-secondary)] mb-[12px]">
      Статус системы
    </div>

    <div className="grid grid-cols-2 gap-[12px]">

      <StatCard
        title="Размещений"
        value={placements.length}
      />

      <StatCard
        title="Бункеров"
        value={bunkers.length}
      />

      <StatCard
        title="Устройства онлайн"
        value={onlineDevices}
        color="var(--status-success-text)"
      />

      <StatCard
        title="Оффлайн"
        value={offlineDevices}
        color="var(--status-danger-text)"
      />

    </div>

  </div>

  {/* еще блок */}

  <div className="bg-[var(--bg-secondary)] rounded-[12px] p-[20px]">

    <div className="text-[14px] text-[var(--text-secondary)] mb-[12px]">
      Производство
    </div>

    <div className="space-y-[12px]">

      <StatCard
        title="Активные размещения"
        value={placements.length}
      />

      {/* <StatCard
        title="Всего рецептов"
        value={recipes.length}
      /> */}

    </div>
  </div>
</div>

  );

};

export default DashboardPage;