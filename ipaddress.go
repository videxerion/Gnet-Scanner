package main

import "net"

// Преобразование IP-адреса в uint32
func ipToUint64(ip net.IP) uint64 {
	return uint64(ip[0])<<24 | uint64(ip[1])<<16 | uint64(ip[2])<<8 | uint64(ip[3])
}

// Преобразование uint32 в IP-адрес
func uint64ToIP(value uint64) net.IP {
	ip := make(net.IP, 4)
	ip[0] = byte(value >> 24)
	ip[1] = byte(value >> 16)
	ip[2] = byte(value >> 8)
	ip[3] = byte(value)
	return ip
}

// nextIP увеличивает IP-адрес на единицу
func nextIP(ip net.IP) net.IP {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			return ip
		}
	}
	return nil // Переполнение
}

func countIPAddresses(cidr string) (uint64, error) {
	// Разбираем CIDR-нотацию
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, err
	}

	// Получаем маску сети и вычисляем количество адресов по маске
	mask := ipnet.Mask
	ones, _ := mask.Size()
	var count uint64 = 1 << (32 - ones) // 2^(32-бит_маски)

	return count, nil
}

func isValidIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return false
	} else if ip.IsMulticast() {
		return false
	} else if ip.IsPrivate() {
		return false
	}

	return true
}
