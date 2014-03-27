<?php
class User_db extends CI_Model {

	private $table	=	'user';

	/**
	 * creates new user in DB
	 *
	 * @author Anil
	 * @param string $email
	 * @param integer $password
	 * @return integer last inserted id
	 */
	function insertNewUser($email,$password){

		$check_duplicate = $this->getUserId($email);

		if($check_duplicate > 0){
			log_message('error','Its a duplicate user with email:'.$email);
			return -1;
		}

		$this->db->insert($this->table, array(
			'email' => $email
			,'password' => $password
		));

		if ($this->db->_error_message()){
			$err_msg = $this->db->_error_message();
			$err_num = $this->db->_error_number();
			log_message('error',$err_num.' : '.$err_msg.'. It happened during inserting new user with email:'.$email.' and password : '.$password);
			return 0;
		}else{
			return $this->db->insert_id();
		}
	}

	function getUserId($email){
		$result = $this->db->get_where($this->table,array('email' => $email));

		if($result->num_rows() > 0){
			//ideally we should get only one row
			$row	= $result->row();
			return $row->user_id;
		}else
			return 0;
	}

	function getUserDetails($email){
		$result = $this->db->get_where($this->table,array('email' => $email));

		if($result->num_rows() > 0){
			//ideally we should get only one row
			return $row	= $result->row();
		}else
			return FALSE;
	}
}