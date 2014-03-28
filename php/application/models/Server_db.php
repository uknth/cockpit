<?php
class Server_db extends CI_Model {

	private $table	=	'server';

	/**
	 * creates new server in DB
	 *
	 * @author Pooja
	 * @param string $name
	 * @param integer $ip
	 * @return integer last inserted id
	 */
	function insertNewServer($name,$ip){

		$check_duplicate = $this->getServerId($name);

		if($check_duplicate > 0){
			log_message('error','Its a duplicate ip:'.$name);
			return -1;
		}

		$this->db->insert($this->table, array(
			'name' => $name
			,'ip' => $ip
		));

		if ($this->db->_error_message()){
			$err_msg = $this->db->_error_message();
			$err_num = $this->db->_error_number();
			log_message('error',$err_num.' : '.$err_msg.'. It happened during inserting new server with name:'.$name.' and ip : '.$ip);
			return 0;
		}else{
			return $this->db->insert_id();
		}
	}

	function getServerId($name){
		$result = $this->db->get_where($this->table,array('name' => $name));

		if($result->num_rows() > 0){
			//ideally we should get only one row
			$row	= $result->row();
			return $row->server_id;
		}else
			return 0;
	}

	function getServerDetails($name){
		$result = $this->db->get_where($this->table,array('name' => $name));

		if($result->num_rows() > 0){
			//ideally we should get only one row
			return $row	= $result->row();
		}else
			return FALSE;
	}
}